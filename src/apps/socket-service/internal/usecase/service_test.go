package usecase

import (
	"context"
	"testing"
	"time"

	"RuntimeRoasters/pkg/base/identity"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestTicketFlow(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	
	service := NewService(rdb, nil, "test-topic", 1*time.Minute, "")
	ctx := context.Background()

	claims := identity.Claims{
		Subject: "user-123",
		Role:    "ADMIN",
	}

	t.Run("Issue and Exchange Ticket", func(t *testing.T) {
		ticket, err := service.IssueTicket(ctx, claims)
		assert.NoError(t, err)
		assert.NotEmpty(t, ticket)

		exchangedClaims, err := service.ExchangeTicket(ctx, ticket)
		assert.NoError(t, err)
		assert.Equal(t, claims.Subject, exchangedClaims.Subject)
		assert.Equal(t, claims.Role, exchangedClaims.Role)

		// One-time use: second exchange should fail
		_, err = service.ExchangeTicket(ctx, ticket)
		assert.Error(t, err)
		assert.Equal(t, "invalid or expired ticket", err.Error())
	})

	t.Run("Invalid Ticket", func(t *testing.T) {
		_, err := service.ExchangeTicket(ctx, "non-existent")
		assert.Error(t, err)
	})
}
