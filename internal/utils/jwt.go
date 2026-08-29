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
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func CreateAccessToken(signingKey []byte, claims JWTClaims, expiresIn int) (string, error) {
	return createToken(signingKey, claims, TokenTypeAccess, expiresIn)
}

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
	// 校验 token 是否为空
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}

	// 使用 jwt 库解析 token
	token, err := jwt.ParseWithClaims(
		tokenString,  // 待解析的 JWT 字符串
		&JWTClaims{}, // 自定义声明结构体，用于存储 token 中的用户信息
		func(token *jwt.Token) (any, error) { // 密钥回调函数，用于验证签名
			// 验证签名算法是否为 HS256（确保 token 使用的是预期的加密算法）
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected jwt signing method: %s", token.Method.Alg())
			}
			// 返回签名密钥，用于验证 token 的签名是否有效
			return signingKey, nil
		},
		// 配置选项：
		// 1. 限制只接受 HS256 签名算法的 token（与上面回调中的检查冗余，但增加了一层安全保障）
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		// 2. 要求 token 必须包含过期时间（exp 声明），如果缺失或已过期则验证失败
		jwt.WithExpirationRequired(),
	)

	// 处理解析错误（包括签名错误、过期错误、格式错误等）
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	return token, nil
}
