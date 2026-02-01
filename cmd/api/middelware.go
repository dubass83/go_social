package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/dubass83/go_social/internal/store"
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

		user, err := app.GetUserFromCacheByID(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		// if it is the user post
		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

		// role precedence check
		allowed, err := app.store.Role.IsPrecedent(r.Context(), user.RoleID, requiredRole)
		if err != nil {
			log.Error().Err(err).Msg("Failed to check role precedence")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if !allowed {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
//      allowed, err := app.store.Role.IsPrecedent(ctx, user.RoleID, roleName)
//      return allowed, err
// }

func (app *application) GetUserFromCacheByID(ctx context.Context, userID int64) (*store.User, error) {
	// Try cache first
	log.Debug().Int64("user_id", userID).Msg("Checking cache for user")
	user, err := app.cache.User.Get(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Int64("user_id", userID).Msg("Failed to get user from cache")
	}
	if user != nil {
		log.Debug().Int64("user", user.ID).Msg("User found in cache")
		return user, nil
	}

	// Cache miss or cache disabled - fetch from database
	log.Debug().
		Int64("user_id", userID).
		Bool("CACHE_ENABLE", app.config.cache.enable).
		Msg("Cache miss or cache disabled - fetching from database")
	user, err = app.store.User.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Store in cache (fire and forget)
	err = app.cache.User.Set(ctx, user)
	if err != nil {
		log.Warn().Err(err).Int64("user_id", userID).Msg("Failed to set user in cache")
	}

	return user, nil
}
