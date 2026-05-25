package token

import (
	"errors"
	"fmt"

	"RuntimeRoasters/pkg/base/auth/provider"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// VerifyAndParseJWT is a framework-agnostic helper to parse and validate JWTs.
func VerifyAndParseJWT(rawToken string, keyProvider provider.KeyProvider, expectedIssuer string) (*identity.Claims, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header[JWTHeaderAlg])
		}
		kid, ok := token.Header[JWTHeaderKeyID].(string)
		if !ok {
			return nil, errors.New("missing kid in token header")
		}

		pubKey, err := keyProvider.GetPublicKey(kid)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch public key: %w", err)
		}

		return pubKey, nil
	}, jwt.WithIssuer(expectedIssuer))

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return extractIdentityClaims(jwtClaims), nil
}

func extractIdentityClaims(jwtClaims jwt.MapClaims) *identity.Claims {
	sub, _ := jwtClaims.GetSubject()
	email := stringClaim(jwtClaims, ClaimEmail)
	orgID := stringClaim(jwtClaims, ClaimOrgID)
	jti, _ := jwtClaims[ClaimJTI].(string)
	storeIDs := stringSliceClaim(jwtClaims, ClaimStoreIDs)
	warehouseIDs := stringSliceClaim(jwtClaims, ClaimWarehouseIDs)

	log := logger.GetLogger().With(zap.String("sub", sub))
	log.Debug("Extracting claims from JWT")

	// Robust Role extraction logic
	var role string
	// 1. Try root level
	if r, ok := jwtClaims[ClaimRole].(string); ok {
		log.Debug("Found role at root level", zap.String("role", r))
		role = r
	} else {
		// 2. Try nested 'ext' level (defensive parsing)
		if extRaw, exists := jwtClaims["ext"]; exists {
			log.Debug("Found 'ext' claim, attempting nested extraction", zap.String("type", fmt.Sprintf("%T", extRaw)))
			switch ext := extRaw.(type) {
			case map[string]interface{}:
				if r, ok := ext[ClaimRole].(string); ok {
					log.Debug("Found role in ext map", zap.String("role", r))
					role = r
				}
			case map[interface{}]interface{}:
				// Handle case where key might be interface{}
				for k, v := range ext {
					if kStr, ok := k.(string); ok && kStr == ClaimRole {
						if vStr, ok := v.(string); ok {
							log.Debug("Found role in interface map", zap.String("role", vStr))
							role = vStr
						}
					}
				}
			default:
				log.Warn("Unexpected 'ext' claim type", zap.String("type", fmt.Sprintf("%T", extRaw)))
			}
		} else {
			log.Debug("No 'role' or 'ext' claim found in JWT")
		}
	}

	return &identity.Claims{
		Subject:      sub,
		Email:        email,
		Role:         role,
		OrgID:        orgID,
		StoreIDs:     storeIDs,
		WarehouseIDs: warehouseIDs,
		JTI:          jti,
	}
}

func stringClaim(claims jwt.MapClaims, key string) string {
	if value, ok := claims[key].(string); ok {
		return value
	}
	if ext, ok := extMap(claims); ok {
		if value, ok := ext[key].(string); ok {
			return value
		}
	}
	return ""
}

func stringSliceClaim(claims jwt.MapClaims, key string) []string {
	if values := toStringSlice(claims[key]); values != nil {
		return values
	}
	if ext, ok := extMap(claims); ok {
		return toStringSlice(ext[key])
	}
	return nil
}

func extMap(claims jwt.MapClaims) (map[string]interface{}, bool) {
	extRaw, exists := claims["ext"]
	if !exists {
		return nil, false
	}
	switch ext := extRaw.(type) {
	case map[string]interface{}:
		return ext, true
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(ext))
		for key, value := range ext {
			if keyString, ok := key.(string); ok {
				out[keyString] = value
			}
		}
		return out, true
	default:
		return nil, false
	}
}

func toStringSlice(raw interface{}) []string {
	switch values := raw.(type) {
	case []string:
		return append([]string(nil), values...)
	case []interface{}:
		out := make([]string, 0, len(values))
		for _, value := range values {
			if str, ok := value.(string); ok && str != "" {
				out = append(out, str)
			}
		}
		return out
	default:
		return nil
	}
}
