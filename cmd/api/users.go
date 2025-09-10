package main

import (
	"net/http"
	"strconv"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
)

type UserPayload struct {
	Username string `json:"username" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (app *application) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		internalServerError(w, r, err)
	}

}

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
