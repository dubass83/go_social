package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != app.config.auth.basic.user || password != app.config.auth.basic.pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) AuthTokenMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Error().Msg("Authorization header is missing")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			log.Error().Msg("Invalid authorization header format")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if parts[0] != "Bearer" {
			log.Error().Msg("Authorization header must start with 'Bearer' but got '" + parts[0] + "'")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			log.Error().Err(err).Msg("Invalid token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !jwtToken.Valid {
			log.Error().Msg("Token is not valid")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, _ := strconv.ParseInt(strconv.Itoa(int(claims["sub"].(float64))), 10, 64)
		ctx := r.Context()

		user, err := app.store.User.GetByID(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
