package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

const secret = "test"

type TestAuthenticator struct{}

var testClaims = jwt.MapClaims{
	"sub": int64(42),
	"iss": "test-token",
	"aud": "test-token",
	"exp": time.Now().Add(time.Hour).Unix(),
}

func (ta *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return "", err
	}
	log.Debug().Msgf("Generated token: %s", tokenString)
	return tokenString, nil
}

func (ta *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
