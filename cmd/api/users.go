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

type UserPayload struct {
	Username string `json:"username" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

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

// CreateUserHandler godoc
//
//	@Summary		Create a new user
//	@Description	create a new user with username, email and password
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UserPayload	true	"User payload"
//	@Success		201		{object}	store.User
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users [post]
func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UserPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err = validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := app.store.User.Create(r.Context(), user); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		internalServerError(w, r, err)
		return
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
