package main

import (
	"net/http"
)

func (app *application) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//Pagination, filter

	ctx := r.Context()

	feed, err := app.store.Post.GetUserFeed(ctx, int64(15))

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		internalServerError(w, r, err)
		return
	}
}
