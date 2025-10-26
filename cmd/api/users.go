package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
)

type UserKey string

const UserCtxKey UserKey = "user"

// GetUserByIDHandler godoc
//
//	@Summary		Get a user
//	@Description	get user by ID
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id} [get]
func (app *application) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		internalServerError(w, r, err)
	}

}

// activateUserHandler godoc
//
//	@Summary		Activate a user
//	@Description	activate user by invitation token
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation Token"
//	@Success		202		{string}	string	"User activated successfully"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put,get]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.User.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			badRequestResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}

	err = app.jsonResponse(w, http.StatusAccepted, "User activated successfully")
	if err != nil {
		internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			badRequestResponse(w, r, err)
			return
		}

		user, err := app.store.User.GetByID(r.Context(), userID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				notFoundResponse(w, r, err)
				return
			default:
				internalServerError(w, r, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user := r.Context().Value(UserCtxKey).(*store.User)
	return user
}
