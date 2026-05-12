package casbin

import (
	"context"
)

type Engine interface {
	Enforce(rvals ...interface{}) (bool, error)
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
