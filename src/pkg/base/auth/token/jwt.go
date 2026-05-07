package token

import (
	"errors"
	"fmt"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/golang-jwt/jwt/v5"
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
	role, _ := jwtClaims[ClaimRole].(string)
	orgID, _ := jwtClaims[ClaimOrgID].(string)
	jti, _ := jwtClaims[ClaimJTI].(string)

	return &identity.Claims{
		Subject: sub,
		Role:    role,
		OrgID:   orgID,
		JTI:     jti,
	}
}
