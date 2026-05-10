package casbin

import (
	"context"

	"github.com/casbin/casbin/v3"
)

type Engine interface {
	Enforce(rvals ...interface{}) (bool, error)
	GetEnforcer() *casbin.Enforcer
	Sync(ctx context.Context) error
}

func MapHTTPMethodToAction(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST", "PUT", "PATCH":
		return "write"
	case "DELETE":
		return "delete"
	default:
		return "read"
	}
}
