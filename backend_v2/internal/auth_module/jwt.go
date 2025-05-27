package authmodule

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secret      []byte
	tokenExpiry time.Duration
}

func NewJWTAuth(secret []byte, expiry time.Duration) *JWTAuth {
	return &JWTAuth{secret: secret, tokenExpiry: expiry}
}

func (m *JWTAuth) GetToken(payload string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": payload,
		"iat":     now.Unix(),
		"exp":     now.Add(m.tokenExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JWTAuth) Validate(tokenStr string) (*string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	}

	id, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid token (user_id)")
	}

	return &id, nil
}
