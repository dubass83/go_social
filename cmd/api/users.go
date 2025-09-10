package main

import (
	"net/http"
	"strconv"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
)

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
