package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"RuntimeRoasters/apps/socket-service/config"
	"RuntimeRoasters/apps/socket-service/internal/usecase"
	"RuntimeRoasters/pkg/base/identity"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRealtimeRoutesRequireTicketForPrivateStream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rdb, cleanup := testRedis(t)
	defer cleanup()

	service := usecase.NewService(rdb, nil, "test-topic", time.Minute, "")
	defer service.Close()
	router := gin.New()
	(&App{Cfg: &config.Config{}, Service: service}).routes(router)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/v1/realtime/stream")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRealtimeRoutesAcceptPublicAndTicketedWebSockets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rdb, cleanup := testRedis(t)
	defer cleanup()

	service := usecase.NewService(rdb, nil, "test-topic", time.Minute, "")
	defer service.Close()
	router := gin.New()
	(&App{Cfg: &config.Config{}, Service: service}).routes(router)
	server := httptest.NewServer(router)
	defer server.Close()

	publicConn, _, err := websocket.DefaultDialer.Dial(wsURL(server.URL, "/v1/realtime/public/topology/ws"), nil)
	require.NoError(t, err)
	require.NoError(t, publicConn.Close())

	ticket, err := service.IssueTicket(context.Background(), identity.Claims{
		Subject:      "warehouse-manager-1",
		Role:         identity.RoleWarehouseMgr,
		WarehouseIDs: []string{"warehouse-1"},
	})
	require.NoError(t, err)

	privateConn, _, err := websocket.DefaultDialer.Dial(wsURL(server.URL, "/v1/realtime/stream?ticket="+ticket), nil)
	require.NoError(t, err)
	require.NoError(t, privateConn.Close())

	_, resp, err := websocket.DefaultDialer.Dial(wsURL(server.URL, "/v1/realtime/stream?ticket="+ticket), nil)
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func testRedis(t *testing.T) (*redis.Client, func()) {
	t.Helper()
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	return client, func() {
		_ = client.Close()
		server.Close()
	}
}

func wsURL(serverURL string, path string) string {
	return "ws" + strings.TrimPrefix(serverURL, "http") + path
}
