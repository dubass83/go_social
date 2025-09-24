package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
)

func (app *application) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	pgFeedQueryDefault := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	pgFeedQuery, err := pgFeedQueryDefault.Parse(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(pgFeedQuery); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Post.GetUserFeed(ctx, int64(97), pgFeedQuery)

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		internalServerError(w, r, err)
		return
	}
}
