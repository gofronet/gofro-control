package jwtutils

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSecretManager struct {
	mu   sync.RWMutex
	priv *rsa.PrivateKey
}

func NewJWTSecretManager() (*JWTSecretManager, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	if err := priv.Validate(); err != nil {
		return nil, err
	}

	return &JWTSecretManager{priv: priv}, nil
}

func (m *JWTSecretManager) Sign(claims jwt.MapClaims) (string, error) {
	m.mu.RLock()
	priv := m.priv
	m.mu.RUnlock()

	if priv == nil {
		return "", fmt.Errorf("jwt private key is not initialized")
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return tok.SignedString(priv)
}

func (m *JWTSecretManager) Verify(tokenString string, expectedIssuer string, expectedAudience string) (jwt.MapClaims, error) {
	m.mu.RLock()
	pub := (*rsa.PublicKey)(nil)
	if m.priv != nil {
		pub = &m.priv.PublicKey
	}
	m.mu.RUnlock()

	if pub == nil {
		return nil, fmt.Errorf("jwt key is not initialized")
	}

	opts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	}

	if expectedIssuer != "" {
		opts = append(opts, jwt.WithIssuer(expectedIssuer))
	}
	if expectedAudience != "" {
		opts = append(opts, jwt.WithAudience(expectedAudience))
	}

	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return pub, nil
	}, opts...)
	if err != nil {
		return nil, err
	}
	if tok == nil || !tok.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return claims, nil
}
