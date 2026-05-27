package usecase

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/topology"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
)

type ClientScope struct {
	Public       bool
	Role         string
	Subject      string
	StoreIDs     []string
	WarehouseIDs []string
	FlowID       string
}

type Client struct {
	id     string
	conn   *websocket.Conn
	send   chan []byte
	scope  ClientScope
	cancel context.CancelFunc
}

type Service struct {
	rdb            *redis.Client
	producer       kafka.Producer
	broadcastTopic string
	sessionTTL     time.Duration
	apiKeys        map[string]apiKey
	clients        map[string]*Client
	mu             sync.RWMutex
	upgrader       websocket.Upgrader
}

type apiKey struct {
	Service string
	Hash    string
	Scopes  map[string]struct{}
}

type Envelope struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewService(rdb *redis.Client, producer kafka.Producer, broadcastTopic string, sessionTTL time.Duration, rawKeys string) *Service {
	if sessionTTL <= 0 {
		sessionTTL = 90 * time.Second
	}
	return &Service{
		rdb:            rdb,
		producer:       producer,
		broadcastTopic: broadcastTopic,
		sessionTTL:     sessionTTL,
		apiKeys:        parseAPIKeys(rawKeys),
		clients:        map[string]*Client{},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *Service) ServeWS(w http.ResponseWriter, r *http.Request, scope ClientScope) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(r.Context())
	client := &Client{id: uuid.NewString(), conn: conn, send: make(chan []byte, 64), scope: scope, cancel: cancel}
	s.addClient(ctx, client)
	go s.writePump(client)
	go s.readPump(client)
	s.sendRecent(ctx, client)
	return nil
}

func (s *Service) IssueTicket(ctx context.Context, claims identity.Claims) (string, error) {
	ticket := uuid.NewString()
	body, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	// Ticket valid for 1 minute
	err = s.rdb.Set(ctx, "socket:ticket:"+ticket, body, 1*time.Minute).Err()
	return ticket, err
}

func (s *Service) ExchangeTicket(ctx context.Context, ticket string) (*identity.Claims, error) {
	if ticket == "" {
		return nil, errors.New("missing ticket")
	}
	key := "socket:ticket:" + ticket
	body, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.New("invalid or expired ticket")
	}
	_ = s.rdb.Del(ctx, key).Err() // One-time use

	var claims identity.Claims
	if err := json.Unmarshal(body, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func (s *Service) HandleKafkaMessage(ctx context.Context, msg kafkago.Message) error {
	cloudEvent, err := events.ParseCloudEvent(msg.Value)
	if err != nil {
		return err
	}
	projection := topology.ProjectEvent(cloudEvent, msg.Value)
	entry := topology.HistoryEntry{
		EventID:        cloudEvent.ID(),
		Topic:          cloudEvent.Type(),
		Status:         stringValue(projection.DisplayPayload, "status"),
		FlowID:         projection.FlowID,
		NodeID:         projection.NodeID,
		EdgeID:         projection.EdgeID,
		Pattern:        projection.Pattern,
		SourceService:  projection.SourceService,
		Visibility:     projection.Visibility,
		TraceID:        projection.TraceID,
		OccurredAt:     cloudEvent.Time(),
		DisplayPayload: projection.DisplayPayload,
	}
	body, err := json.Marshal(Envelope{Type: "topology.event", Data: entry})
	if err != nil {
		return err
	}
	_ = s.rdb.Set(ctx, "socket:last-event:"+entry.FlowID, body, s.sessionTTL).Err()
	s.broadcast(entry, body)
	return nil
}

func (s *Service) PublishInternalEvent(ctx context.Context, keyHeader string, event topology.BroadcastRequested) error {
	if err := s.authorizeAPIKey(keyHeader, "events:publish"); err != nil {
		return err
	}
	if event.EventID == "" {
		event.EventID = uuid.NewString()
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}
	if event.Visibility == "" {
		event.Visibility = topology.VisibilityPrivate
	}
	payload := events.SocketBroadcastRequested{
		EventID:       event.EventID,
		Channel:       event.Channel,
		FlowID:        event.FlowID,
		NodeID:        event.NodeID,
		EdgeID:        event.EdgeID,
		Visibility:    event.Visibility,
		Status:        event.Status,
		SourceService: event.SourceService,
		EventType:     "socket.broadcast.requested",
		Payload:       event.Payload,
		OccurredAt:    event.OccurredAt,
	}
	cloudEvent, err := events.NewCloudEvent(ctx, s.broadcastTopic, "/services/socket-service", "socket/broadcasts/"+event.EventID, payload, events.Metadata{
		EventID:       event.EventID,
		CorrelationID: event.EventID,
		OccurredAt:    event.OccurredAt,
	})
	if err != nil {
		return err
	}
	return s.producer.Publish(ctx, s.broadcastTopic, event.EventID, cloudEvent)
}

func (s *Service) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, client := range s.clients {
		client.cancel()
		_ = client.conn.Close()
		close(client.send)
	}
	s.clients = map[string]*Client{}
}

func (s *Service) addClient(ctx context.Context, client *Client) {
	s.mu.Lock()
	s.clients[client.id] = client
	s.mu.Unlock()
	value := map[string]interface{}{"role": client.scope.Role, "flow_id": client.scope.FlowID, "public": client.scope.Public}
	body, _ := json.Marshal(value)
	_ = s.rdb.Set(ctx, "socket:session:"+client.id, body, s.sessionTTL).Err()
	_ = s.rdb.Incr(ctx, "socket:presence:"+presenceKey(client.scope)).Err()
	_ = s.rdb.Expire(ctx, "socket:presence:"+presenceKey(client.scope), s.sessionTTL).Err()
}

func (s *Service) removeClient(client *Client) {
	s.mu.Lock()
	if _, ok := s.clients[client.id]; ok {
		delete(s.clients, client.id)
		close(client.send)
	}
	s.mu.Unlock()
	_ = s.rdb.Del(context.Background(), "socket:session:"+client.id).Err()
	_ = s.rdb.Decr(context.Background(), "socket:presence:"+presenceKey(client.scope)).Err()
}

func (s *Service) broadcast(entry topology.HistoryEntry, body []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, client := range s.clients {
		if !canReceive(client.scope, entry) {
			continue
		}
		select {
		case client.send <- body:
		default:
		}
	}
}

func (s *Service) sendRecent(ctx context.Context, client *Client) {
	if client.scope.FlowID == "" {
		return
	}
	if body, err := s.rdb.Get(ctx, "socket:last-event:"+client.scope.FlowID).Bytes(); err == nil {
		client.send <- body
	}
}

func (s *Service) writePump(client *Client) {
	defer client.cancel()
	for body := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, body); err != nil {
			return
		}
	}
}

func (s *Service) readPump(client *Client) {
	defer func() {
		client.cancel()
		_ = client.conn.Close()
		s.removeClient(client)
	}()
	for {
		if _, _, err := client.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (s *Service) authorizeAPIKey(raw string, scope string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return errors.New("missing api key")
	}
	sum := sha256.Sum256([]byte(raw))
	hash := hex.EncodeToString(sum[:])
	for _, key := range s.apiKeys {
		if subtle.ConstantTimeCompare([]byte(hash), []byte(key.Hash)) != 1 {
			continue
		}
		if _, ok := key.Scopes[scope]; !ok {
			return fmt.Errorf("missing scope %s", scope)
		}
		return nil
	}
	return errors.New("invalid api key")
}

func parseAPIKeys(raw string) map[string]apiKey {
	keys := map[string]apiKey{}
	for _, part := range strings.Split(raw, ",") {
		fields := strings.SplitN(strings.TrimSpace(part), ":", 3)
		if len(fields) != 3 {
			continue
		}
		scopes := map[string]struct{}{}
		for _, scope := range strings.Split(fields[2], "|") {
			scope = strings.TrimSpace(scope)
			if scope != "" {
				scopes[scope] = struct{}{}
			}
		}
		keys[fields[0]] = apiKey{Service: fields[0], Hash: fields[1], Scopes: scopes}
	}
	return keys
}

func canReceive(scope ClientScope, entry topology.HistoryEntry) bool {
	if scope.FlowID != "" && entry.FlowID != scope.FlowID {
		return false
	}
	if scope.Public {
		return entry.Visibility == topology.VisibilityPublic
	}
	if scope.Role == identity.RoleAdmin {
		return true
	}
	storeID := stringValue(entry.DisplayPayload, "store_id")
	if scope.Role == identity.RoleStoreMgr && storeID != "" {
		return contains(scope.StoreIDs, storeID)
	}
	warehouseID := stringValue(entry.DisplayPayload, "warehouse_id")
	if scope.Role == identity.RoleWarehouseMgr && warehouseID != "" {
		return contains(scope.WarehouseIDs, warehouseID)
	}
	driverID := stringValue(entry.DisplayPayload, "driver_id")
	if scope.Role == "DRIVER" && driverID != "" {
		return driverID == scope.Subject
	}
	return entry.Visibility == topology.VisibilityPublic
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func stringValue(payload map[string]interface{}, key string) string {
	value, _ := payload[key].(string)
	return value
}

func presenceKey(scope ClientScope) string {
	if scope.Public {
		return "public"
	}
	if scope.Role != "" {
		return scope.Role
	}
	return "unknown"
}
