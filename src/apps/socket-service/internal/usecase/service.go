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

type Notification struct {
	ID             string     `json:"id"`
	TargetRole     string     `json:"target_role,omitempty"`
	TargetUserID   string     `json:"target_user_id,omitempty"`
	FarmID         string     `json:"farm_id,omitempty"`
	WarehouseID    string     `json:"warehouse_id,omitempty"`
	StoreID        string     `json:"store_id,omitempty"`
	ShipmentID     string     `json:"shipment_id,omitempty"`
	EntityType     string     `json:"entity_type"`
	EntityID       string     `json:"entity_id"`
	Type           string     `json:"type"`
	Title          string     `json:"title"`
	Message        string     `json:"message"`
	Severity       string     `json:"severity"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
}

const (
	NotificationStatusUnread   = "UNREAD"
	NotificationStatusAcked    = "ACKED"
	NotificationStatusResolved = "RESOLVED"
)

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
	if notification, ok := notificationFromEvent(entry, cloudEvent); ok {
		if err := s.CreateNotification(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) CreateNotification(ctx context.Context, notification Notification) error {
	if notification.ID == "" {
		notification.ID = uuid.NewString()
	}
	if notification.Status == "" {
		notification.Status = NotificationStatusUnread
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	key := "notification:" + notification.ID
	if err := s.rdb.Set(ctx, key, body, 14*24*time.Hour).Err(); err != nil {
		return err
	}
	for _, index := range notificationIndexes(notification) {
		if err := s.rdb.LPush(ctx, index, notification.ID).Err(); err != nil {
			return err
		}
		_ = s.rdb.LTrim(ctx, index, 0, 199).Err()
		_ = s.rdb.Expire(ctx, index, 14*24*time.Hour).Err()
	}
	notificationBody, err := json.Marshal(Envelope{Type: "notification.created", Data: notification})
	if err != nil {
		return err
	}
	s.broadcastNotification(notification, notificationBody)
	return nil
}

func (s *Service) ListNotifications(ctx context.Context, claims identity.Claims) ([]Notification, error) {
	ids := make([]string, 0, 128)
	seen := map[string]struct{}{}
	for _, index := range notificationIndexesForClaims(claims) {
		values, err := s.rdb.LRange(ctx, index, 0, 99).Result()
		if err != nil {
			return nil, err
		}
		for _, id := range values {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	notifications := make([]Notification, 0, len(ids))
	for _, id := range ids {
		body, err := s.rdb.Get(ctx, "notification:"+id).Bytes()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		var notification Notification
		if err := json.Unmarshal(body, &notification); err != nil {
			return nil, err
		}
		if canReceiveNotification(claims, notification) {
			notifications = append(notifications, notification)
		}
	}
	return notifications, nil
}

func (s *Service) AcknowledgeNotification(ctx context.Context, claims identity.Claims, id string) (*Notification, error) {
	return s.updateNotificationStatus(ctx, claims, id, NotificationStatusAcked)
}

func (s *Service) ResolveNotification(ctx context.Context, claims identity.Claims, id string) (*Notification, error) {
	return s.updateNotificationStatus(ctx, claims, id, NotificationStatusResolved)
}

func (s *Service) updateNotificationStatus(ctx context.Context, claims identity.Claims, id string, status string) (*Notification, error) {
	key := "notification:" + id
	body, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var notification Notification
	if err := json.Unmarshal(body, &notification); err != nil {
		return nil, err
	}
	if !canReceiveNotification(claims, notification) {
		return nil, errors.New("notification access denied")
	}
	now := time.Now()
	notification.Status = status
	if status == NotificationStatusAcked {
		notification.AcknowledgedAt = &now
	}
	updated, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}
	if err := s.rdb.Set(ctx, key, updated, 14*24*time.Hour).Err(); err != nil {
		return nil, err
	}
	return &notification, nil
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

func (s *Service) broadcastNotification(notification Notification, body []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, client := range s.clients {
		if client.scope.Public || !canReceiveNotification(scopeToClaims(client.scope), notification) {
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

func notificationFromEvent(entry topology.HistoryEntry, event interface {
	Extensions() map[string]interface{}
}) (Notification, bool) {
	notification := Notification{
		ID:          "n-" + entry.EventID,
		EntityType:  entityType(entry.Topic),
		EntityID:    firstNonEmpty(extensionValue(event, "orderid"), extensionValue(event, "shipmentid"), extensionValue(event, "harvestid"), entry.EventID),
		Type:        entry.Topic,
		Severity:    "info",
		Status:      NotificationStatusUnread,
		CreatedAt:   entry.OccurredAt,
		StoreID:     extensionValue(event, "storeid"),
		WarehouseID: extensionValue(event, "warehouseid"),
		FarmID:      extensionValue(event, "farmid"),
		ShipmentID:  extensionValue(event, "shipmentid"),
	}
	switch entry.Topic {
	case events.TopicFarmHarvestCreated, events.TopicWarehousePickupRequested:
		notification.TargetRole = identity.RoleWarehouseMgr
		notification.Title = "Pickup work available"
		notification.Message = "A harvest is ready for warehouse pickup coordination."
	case events.TopicWarehouseDispatchRequested:
		notification.TargetRole = identity.RoleWarehouseMgr
		notification.Title = "Outbound dispatch requested"
		notification.Message = "A paid order has reserved stock and is ready for delivery dispatch."
	case events.TopicLogisticsDeliveryAssigned:
		notification.TargetRole = "DRIVER"
		notification.TargetUserID = extensionValue(event, "driverid")
		notification.Title = "Delivery assigned"
		notification.Message = "A retail delivery has been assigned to a driver."
	case events.TopicLogisticsDeliveryDriverConfirmed, events.TopicLogisticsDriverReturnedToBase:
		notification.TargetRole = identity.RoleStoreMgr
		notification.Title = "Delivery status updated"
		notification.Message = "A retail delivery status changed."
	case events.TopicWarehouseStockReservationFailed, events.TopicPaymentFailed:
		notification.TargetRole = identity.RoleStoreMgr
		notification.Severity = "error"
		notification.Title = "Order fulfillment failed"
		notification.Message = "A paid order could not continue through fulfillment."
	default:
		return Notification{}, false
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}
	return notification, true
}

func extensionValue(event interface {
	Extensions() map[string]interface{}
}, key string) string {
	value, ok := event.Extensions()[key]
	if !ok || value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func notificationIndexes(notification Notification) []string {
	indexes := []string{}
	if notification.TargetRole != "" {
		indexes = append(indexes, "notifications:role:"+notification.TargetRole)
	}
	if notification.TargetUserID != "" {
		indexes = append(indexes, "notifications:user:"+notification.TargetUserID)
	}
	if notification.StoreID != "" {
		indexes = append(indexes, "notifications:store:"+notification.StoreID)
	}
	if notification.WarehouseID != "" {
		indexes = append(indexes, "notifications:warehouse:"+notification.WarehouseID)
	}
	return indexes
}

func notificationIndexesForClaims(claims identity.Claims) []string {
	indexes := []string{"notifications:role:" + claims.Role, "notifications:user:" + claims.Subject}
	for _, storeID := range claims.StoreIDs {
		indexes = append(indexes, "notifications:store:"+storeID)
	}
	for _, warehouseID := range claims.WarehouseIDs {
		indexes = append(indexes, "notifications:warehouse:"+warehouseID)
	}
	return indexes
}

func canReceiveNotification(claims identity.Claims, notification Notification) bool {
	if claims.IsAdmin() {
		return true
	}
	if notification.TargetUserID != "" && notification.TargetUserID != claims.Subject {
		return false
	}
	if notification.TargetRole != "" && notification.TargetRole != claims.Role {
		return false
	}
	if claims.Role == identity.RoleStoreMgr && notification.StoreID != "" && !claims.CanAccessStore(notification.StoreID) {
		return false
	}
	if claims.Role == identity.RoleWarehouseMgr && notification.WarehouseID != "" && !claims.CanAccessWarehouse(notification.WarehouseID) {
		return false
	}
	return notification.TargetRole != "" || notification.TargetUserID != "" || notification.StoreID != "" || notification.WarehouseID != ""
}

func scopeToClaims(scope ClientScope) identity.Claims {
	return identity.Claims{
		Subject:      scope.Subject,
		Role:         scope.Role,
		StoreIDs:     scope.StoreIDs,
		WarehouseIDs: scope.WarehouseIDs,
	}
}

func entityType(topic string) string {
	switch {
	case strings.Contains(topic, "order"), strings.Contains(topic, "payment"), strings.Contains(topic, "delivery"):
		return "order"
	case strings.Contains(topic, "pickup"), strings.Contains(topic, "harvest"):
		return "harvest"
	default:
		return "event"
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
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
