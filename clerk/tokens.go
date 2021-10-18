package clerk

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"

	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	standardClaimsKeys = []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"}
)

type TokenClaims struct {
	jwt.Claims
	Extra map[string]interface{}
}

type SessionClaims struct {
	jwt.Claims
	SessionID string `json:"sid"`
}

// DecodeToken decodes a jwt token without verifying it.
func (c *client) DecodeToken(token string) (*TokenClaims, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	standardClaims := jwt.Claims{}
	extraClaims := make(map[string]interface{})

	if err = parsedToken.UnsafeClaimsWithoutVerification(&standardClaims, &extraClaims); err != nil {
		return nil, err
	}

	// Delete any standard claims included in the extra claims
	for _, key := range standardClaimsKeys {
		delete(extraClaims, key)
	}

	return &TokenClaims{Claims: standardClaims, Extra: extraClaims}, nil
}

// VerifyToken verifies the session jwt token.
func (c *client) VerifyToken(token string) (*SessionClaims, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	if len(parsedToken.Headers) == 0 {
		return nil, fmt.Errorf("missing jwt headers")
	}

	kid := parsedToken.Headers[0].KeyID
	if kid == "" {
		return nil, fmt.Errorf("missing jwt kid header claim")
	}

	jwk, err := c.getJWK(kid)
	if err != nil {
		return nil, err
	}

	if parsedToken.Headers[0].Algorithm != jwk.Algorithm {
		return nil, fmt.Errorf("invalid signing algorithm %s", jwk.Algorithm)
	}

	claims := SessionClaims{}
	if err = parsedToken.Claims(jwk.Key, &claims); err != nil {
		return nil, err
	}

	if err = claims.Claims.ValidateWithLeeway(jwt.Expected{Time: time.Now()}, 0); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(claims.Issuer, "https://clerk.") {
		return nil, fmt.Errorf("invalid issuer %s", claims.Issuer)
	}

	return &claims, nil
}

func (c *client) getJWK(kid string) (*jose.JSONWebKey, error) {
	if c.jwksCache.isInvalid() {
		jwks, err := c.jwks.ListAll()
		if err != nil {
			return nil, err
		}

		c.jwksCache.set(jwks)
	}

	return c.jwksCache.get(kid)
}