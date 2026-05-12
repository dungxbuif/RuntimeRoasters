package token

import (
	"errors"
	"fmt"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
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
	orgID, _ := jwtClaims[ClaimOrgID].(string)
	jti, _ := jwtClaims[ClaimJTI].(string)

	log := logger.GetLogger().With(zap.String("sub", sub))
	
	fmt.Printf("[DEBUG] JWT: Extracting claims for sub: %s\n", sub)
	log.Debug("Extracting claims from JWT")

	// Robust Role extraction logic
	var role string
	// 1. Try root level
	if r, ok := jwtClaims[ClaimRole].(string); ok {
		fmt.Printf("[DEBUG] JWT: Found role at root level: %s\n", r)
		log.Debug("Found role at root level", zap.String("role", r))
		role = r
	} else {
		// 2. Try nested 'ext' level (defensive parsing)
		if extRaw, exists := jwtClaims["ext"]; exists {
			fmt.Printf("[DEBUG] JWT: Found 'ext' claim. Type: %T\n", extRaw)
			log.Debug("Found 'ext' claim, attempting nested extraction", zap.String("type", fmt.Sprintf("%T", extRaw)))
			switch ext := extRaw.(type) {
			case map[string]interface{}:
				if r, ok := ext[ClaimRole].(string); ok {
					fmt.Printf("[DEBUG] JWT: Found role in ext map: %s\n", r)
					log.Debug("Found role in ext map", zap.String("role", r))
					role = r
				}
			case map[interface{}]interface{}:
				// Handle case where key might be interface{}
				for k, v := range ext {
					if kStr, ok := k.(string); ok && kStr == ClaimRole {
						if vStr, ok := v.(string); ok {
							fmt.Printf("[DEBUG] JWT: Found role in interface map: %s\n", vStr)
							log.Debug("Found role in interface map", zap.String("role", vStr))
							role = vStr
						}
					}
				}
			default:
				fmt.Printf("[DEBUG] JWT: 'ext' claim is unexpected type: %T\n", extRaw)
				log.Warn("Unexpected 'ext' claim type", zap.String("type", fmt.Sprintf("%T", extRaw)))
			}
		} else {
			fmt.Printf("[DEBUG] JWT: No 'role' or 'ext' claim found\n")
			log.Debug("No 'role' or 'ext' claim found in JWT")
		}
	}

	return &identity.Claims{
		Subject: sub,
		Role:    role,
		OrgID:   orgID,
		JTI:     jti,
	}
}
