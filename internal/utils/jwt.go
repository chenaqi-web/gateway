package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type JWTClaims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role,omitempty"`
	// 预留给Core校验token版本
	AuthVersion uint64 `json:"auth_version,omitempty"`
	TokenType   string `json:"token_type"`
	jwt.RegisteredClaims
}

// 创建access token
func CreateAccessToken(signingKey []byte, claims JWTClaims, expiresIn int) (string, error) {
	return createToken(signingKey, claims, TokenTypeAccess, expiresIn)
}

// 创建refresh token
func CreateRefreshToken(signingKey []byte, claims JWTClaims, expiresIn int) (string, error) {
	return createToken(signingKey, claims, TokenTypeRefresh, expiresIn)
}

func createToken(signingKey []byte, claims JWTClaims, tokenType string, expiresIn int) (string, error) {
	now := time.Now().UTC()
	claims.TokenType = tokenType
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(expiresIn) * time.Second))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signedToken, nil
}

// 获取token信息
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

// 解析token
func ParseToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
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
