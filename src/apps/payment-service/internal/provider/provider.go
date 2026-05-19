package provider

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ProviderStripe = "STRIPE"
	ProviderVNPay  = "VNPAY"
)

type IntentRequest struct {
	PaymentID string
	OrderID   string
	Amount    float64
	Currency  string
}

type IntentResponse struct {
	ProviderRef string
	CheckoutURL string
	Simulated   bool
}

type RefundResponse struct {
	RefundRef string
	Simulated bool
}

type Gateway interface {
	Name() string
	CreateIntent(ctx context.Context, req IntentRequest) (IntentResponse, error)
	Refund(ctx context.Context, providerRef string, amount float64, currency string, reason string) (RefundResponse, error)
	VerifyWebhook(payload []byte, timestamp string, signature string) error
}

type Factory struct {
	gateways map[string]Gateway
}

func NewFactory(stripeSecret string, vnpaySecret string) *Factory {
	return &Factory{gateways: map[string]Gateway{
		ProviderStripe: newSimulatedGateway(ProviderStripe, "pi_sim", "re_sim", "stripe", stripeSecret, sha256Hash),
		ProviderVNPay:  newSimulatedGateway(ProviderVNPay, "vnpay_sim", "vnpay_refund", "vnpay", vnpaySecret, sha512Hash),
	}}
}

func (f *Factory) Get(name string) (Gateway, error) {
	key := strings.ToUpper(strings.TrimSpace(name))
	if key == "" {
		key = ProviderStripe
	}
	gateway, ok := f.gateways[key]
	if !ok {
		return nil, fmt.Errorf("unsupported payment provider %s", name)
	}
	return gateway, nil
}

type simulatedGateway struct {
	name          string
	intentPrefix  string
	refundPrefix  string
	checkoutHost  string
	webhookSecret string
	hash          func(secret string, payload []byte) string
}

func newSimulatedGateway(name string, intentPrefix string, refundPrefix string, checkoutHost string, webhookSecret string, hash func(string, []byte) string) Gateway {
	return &simulatedGateway{name: name, intentPrefix: intentPrefix, refundPrefix: refundPrefix, checkoutHost: checkoutHost, webhookSecret: webhookSecret, hash: hash}
}

func (g *simulatedGateway) Name() string {
	return g.name
}

func (g *simulatedGateway) CreateIntent(ctx context.Context, req IntentRequest) (IntentResponse, error) {
	_ = ctx
	ref := fmt.Sprintf("%s_%s", g.intentPrefix, uuid.NewString())
	return IntentResponse{
		ProviderRef: ref,
		CheckoutURL: fmt.Sprintf("https://simulate.%s.local/checkout/%s?order_id=%s", g.checkoutHost, ref, req.OrderID),
		Simulated:   true,
	}, nil
}

func (g *simulatedGateway) Refund(ctx context.Context, providerRef string, amount float64, currency string, reason string) (RefundResponse, error) {
	_ = ctx
	_ = amount
	_ = currency
	_ = reason
	return RefundResponse{RefundRef: fmt.Sprintf("%s_%s_%s", g.refundPrefix, providerRef, uuid.NewString()), Simulated: true}, nil
}

func (g *simulatedGateway) VerifyWebhook(payload []byte, timestamp string, signature string) error {
	if g.webhookSecret == "" {
		return errors.New("webhook secret is not configured")
	}
	if timestamp == "" || signature == "" {
		return errors.New("webhook timestamp and signature are required")
	}
	parsed, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return errors.New("invalid webhook timestamp")
	}
	signedAt := time.Unix(parsed, 0)
	now := time.Now()
	if signedAt.Before(now.Add(-5*time.Minute)) || signedAt.After(now.Add(5*time.Minute)) {
		return errors.New("webhook timestamp outside tolerance")
	}
	expected := g.hash(g.webhookSecret, []byte(timestamp+"."+string(payload)))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return errors.New("invalid webhook signature")
	}
	return nil
}

func StripeSignatureTimestamp(header string) (string, string) {
	var timestamp string
	var signature string
	for _, part := range strings.Split(header, ",") {
		keyValue := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(keyValue) != 2 {
			continue
		}
		switch keyValue[0] {
		case "t":
			timestamp = keyValue[1]
		case "v1":
			signature = keyValue[1]
		}
	}
	return timestamp, signature
}

func sha256Hash(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func sha512Hash(secret string, payload []byte) string {
	mac := hmac.New(sha512.New, []byte(secret))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
