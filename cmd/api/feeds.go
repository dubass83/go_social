package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
)

// GetUserFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	get paginated feed of posts from followed users
//	@Tags			FEEDS
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit number of posts"	default(10)
//	@Param			offset	query		int		false	"Offset for pagination"	default(0)
//	@Param			sort	query		string	false	"Sort order (asc/desc)"	default(desc)
//	@Param			tags	query		string	false	"Filter by tags (comma-separated)"
//	@Param			search	query		string	false	"Search query"
//	@Success		200		{array}		store.Post
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users/feed [get]
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

	user := getUserFromCtx(r)

	ctx := r.Context()

	feed, err := app.store.Post.GetUserFeed(ctx, user.ID, pgFeedQuery)

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		internalServerError(w, r, err)
		return
	}
}
