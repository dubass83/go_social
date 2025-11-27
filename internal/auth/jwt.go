package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

func (jwta *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwta.iss
	token.Header["aud"] = jwta.aud
	token.Header["iss"] = jwta.iss
	return token.SignedString([]byte(jwta.secret))
}

func (jwta *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwta.secret), nil
	})
}
