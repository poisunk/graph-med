package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID   string
	Username string
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 接受者
			Audience: jwt.ClaimStrings{username},
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			// 唯一表示符
			ID: userID,
			// 签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 签发者
			Issuer: "graph-med",
			// 生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			// 主题
			Subject: "token",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString([]byte(JWT_SECRET)); err != nil {
		return "", err
	} else {
		return token, nil
	}
}

// RefreshToken 刷新JWT Token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return GenerateToken(claims.UserID, claims.Username)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析JWT Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(JWT_SECRET), nil
	})

	// 如果解析成功
	if err == nil && token != nil {
		claim, ok := token.Claims.(*Claims)
		if ok && token.Valid {
			return claim, nil
		}
	}
	return nil, err
}
