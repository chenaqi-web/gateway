package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess       = "access"
	TokenTypeRefresh      = "refresh"
	minJWTSigningKeyBytes = 32
)

type JWTClaims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role,omitempty"`
	// Reserved for future token revocation checks against Core user auth versions.
	AuthVersion uint64 `json:"auth_version,omitempty"`
	TokenType   string `json:"token_type"`
	jwt.RegisteredClaims
}

func CreateAccessToken(signingKey []byte, claims JWTClaims, expiresIn time.Duration) (string, error) {
	return createToken(signingKey, claims, TokenTypeAccess, expiresIn)
}

func CreateRefreshToken(signingKey []byte, claims JWTClaims, expiresIn time.Duration) (string, error) {
	return createToken(signingKey, claims, TokenTypeRefresh, expiresIn)
}

func createToken(signingKey []byte, claims JWTClaims, tokenType string, expiresIn time.Duration) (string, error) {
	if len(signingKey) < minJWTSigningKeyBytes {
		return "", fmt.Errorf("jwt signing key must be at least %d bytes", minJWTSigningKeyBytes)
	}
	if expiresIn <= 0 {
		return "", fmt.Errorf("token expiration must be positive")
	}

	now := time.Now().UTC()
	claims.TokenType = tokenType
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(expiresIn))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signedToken, nil
}

// =====================================================================================================================

func GetClaims(tokenString string, signingKey []byte) (*JWTClaims, error) {
	token, err := ParseToken(tokenString, signingKey)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token claims are invalid")
	}
	return claims, nil
}

func ParseToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}
	if len(signingKey) < minJWTSigningKeyBytes {
		return nil, fmt.Errorf("jwt signing key must be at least %d bytes", minJWTSigningKeyBytes)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected jwt signing method: %s", token.Method.Alg())
			}
			return signingKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	return token, nil
}
