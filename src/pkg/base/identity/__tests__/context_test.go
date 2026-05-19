package identity_tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
)

func TestInjectAndFromContext(t *testing.T) {
	t.Run("happy path: inject then retrieve", func(t *testing.T) {
		want := identity.Claims{
			Subject: "user-123",
			Role:    "admin",
			OrgID:   "org-abc",
			JTI:     "jti-xyz",
		}

		ctx := identity.InjectContext(context.Background(), want)
		got, ok := identity.FromContext(ctx)

		if !ok {
			t.Fatal("expected identity in context, got none")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("empty context returns false", func(t *testing.T) {
		_, ok := identity.FromContext(context.Background())
		if ok {
			t.Error("expected no identity in empty context")
		}
	})
}

func TestMustFromContext_Panics(t *testing.T) {
	t.Run("panics when no identity in context", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic but did not get one")
			}
		}()
		identity.MustFromContext(context.Background())
	})

	t.Run("does not panic when identity is present", func(t *testing.T) {
		c := identity.Claims{Subject: "user-456"}
		ctx := identity.InjectContext(context.Background(), c)
		got := identity.MustFromContext(ctx)
		if got.Subject != "user-456" {
			t.Errorf("got %s, want user-456", got.Subject)
		}
	})
}
