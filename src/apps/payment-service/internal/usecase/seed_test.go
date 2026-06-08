package usecase_test

import (
	"context"
	"testing"

	paymentdomain "RuntimeRoasters/apps/payment-service/internal/domain"
	"RuntimeRoasters/apps/payment-service/internal/provider"
	paymentusecase "RuntimeRoasters/apps/payment-service/internal/usecase"
	"RuntimeRoasters/pkg/events"
	"github.com/stretchr/testify/require"
)

func TestSeedDemoDataCreatesDeterministicFinanceCases(t *testing.T) {
	ctx := context.Background()
	db := openPaymentDB(t)
	producer := &recordingProducer{}
	service := paymentusecase.NewService(
		db,
		producer,
		provider.NewFactory("stripe-secret", "vnpay-secret"),
		provider.ProviderStripe,
		"VND",
		true,
		events.TopicPaymentIntentCreated,
		events.TopicPaymentSimulatedCompleted,
		events.TopicPaymentFailed,
		events.TopicPaymentRefunded,
	)

	created, err := service.SeedDemoData(ctx, false, map[string]string{})
	require.NoError(t, err)
	require.Equal(t, int64(3), created)

	createdAgain, err := service.SeedDemoData(ctx, false, map[string]string{})
	require.NoError(t, err)
	require.Equal(t, int64(3), createdAgain)

	count, err := service.DemoPaymentCount(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(3), count)

	var payments []paymentdomain.Payment
	require.NoError(t, db.Order("provider_ref ASC").Find(&payments).Error)
	require.Len(t, payments, 3)
	require.Contains(t, statuses(payments), paymentdomain.PaymentStatusSucceeded)
	require.Contains(t, statuses(payments), paymentdomain.PaymentStatusFailed)
	require.Contains(t, statuses(payments), paymentdomain.PaymentStatusRefunded)
	require.NotEmpty(t, payments[2].Items)
	require.NotEmpty(t, producer.events)
}

func statuses(payments []paymentdomain.Payment) []paymentdomain.PaymentStatus {
	out := make([]paymentdomain.PaymentStatus, 0, len(payments))
	for _, payment := range payments {
		out = append(out, payment.Status)
	}
	return out
}
