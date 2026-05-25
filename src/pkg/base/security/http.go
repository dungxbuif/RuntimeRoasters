package security

import (
	"context"
	"time"

	"RuntimeRoasters/pkg/base/auth/provider"
	authhttp "RuntimeRoasters/pkg/base/auth/transport/http"
	"RuntimeRoasters/pkg/base/casbin"
	casbingrpc "RuntimeRoasters/pkg/base/casbin/transport/grpc"
	casbinhttp "RuntimeRoasters/pkg/base/casbin/transport/http"
	"RuntimeRoasters/pkg/kafka"
	"github.com/gin-gonic/gin"
)

const DefaultRBACModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act))
`

type HTTPOptions struct {
	ServiceName      string
	JWKSURL          string
	InternalSecret   string
	JWKSCacheTTL     string
	ExpectedIssuer   string
	AuthServiceAddr  string
	KafkaBrokers     []string
	KafkaPolicyTopic string
}

func (o *HTTPOptions) ApplyDefaults() {
	if o.InternalSecret == "" {
		o.InternalSecret = "68f59c82c34424dbd0853a6cf159047be593d4e98d766db7e5047ec9ead3c71c"
	}
	if o.JWKSURL == "" {
		o.JWKSURL = "http://localhost:4434/.well-known/jwks.json"
	}
	if o.JWKSCacheTTL == "" {
		o.JWKSCacheTTL = "5m"
	}
	if o.ExpectedIssuer == "" {
		o.ExpectedIssuer = "http://localhost:4444/"
	}
	if o.AuthServiceAddr == "" {
		o.AuthServiceAddr = "localhost:50052"
	}
	if o.KafkaPolicyTopic == "" {
		o.KafkaPolicyTopic = "auth.policy.changed"
	}
}

type HTTPGuards struct {
	KeyProvider    provider.KeyProvider
	Enforcer       *casbin.ResilientReader
	PolicyConsumer kafka.Consumer
	Authn          gin.HandlerFunc
	Authz          gin.HandlerFunc
}

func NewHTTPGuards(opts HTTPOptions) (*HTTPGuards, error) {
	opts.ApplyDefaults()
	ttl, _ := time.ParseDuration(opts.JWKSCacheTTL)
	keyProvider, err := provider.NewJWKSCache(opts.JWKSURL, opts.InternalSecret, ttl)
	if err != nil {
		return nil, err
	}

	policyConsumer := kafka.NewConsumer(opts.KafkaBrokers, opts.ServiceName+"-auth-policy", opts.KafkaPolicyTopic)

	authSnapshotClient, err := casbingrpc.NewAuthSnapshotClient(opts.AuthServiceAddr, opts.ServiceName)
	if err != nil {
		_ = policyConsumer.Close()
		return nil, err
	}
	reader, err := casbin.NewResilientReader(authSnapshotClient, casbin.ReaderOptions{
		ModelText:            DefaultRBACModel,
		SyncInterval:         5 * time.Minute,
		RetryInterval:        2 * time.Second,
		PolicyChangeConsumer: policyConsumer,
	})
	if err != nil {
		_ = policyConsumer.Close()
		return nil, err
	}
	reader.StartBackgroundSync(context.Background())

	return &HTTPGuards{
		KeyProvider:    keyProvider,
		Enforcer:       reader,
		PolicyConsumer: policyConsumer,
		Authn:          authhttp.GinMiddleware(keyProvider, opts.ExpectedIssuer),
		Authz:          casbinhttp.GinMiddleware(reader),
	}, nil
}

func (g *HTTPGuards) Close() {
	if g == nil || g.PolicyConsumer == nil {
		return
	}
	_ = g.PolicyConsumer.Close()
}
