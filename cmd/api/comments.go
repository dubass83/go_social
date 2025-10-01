package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
)

type commentPayLoad struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// CreateCommentToPostByIDHandler godoc
//
//	@Summary		Create a comment on a post
//	@Description	create a new comment on a post by post ID
//	@Tags			COMMENTS
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Post ID"
//	@Param			comment	body		commentPayLoad	true	"Comment payload"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/posts/{id}/comments [post]
func (app *application) CreateCommentToPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getPostFromCtx(r)

	payload := commentPayLoad{}

	if err := readJSON(w, r, &payload); err != nil {
		internalServerError(w, r, err)
	}

	if err := validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		UserID:  payload.UserID,
		Content: payload.Content,
		PostID:  post.ID,
	}

	if err := app.store.Comment.Create(ctx, comment); err != nil {
		internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		internalServerError(w, r, err)
		return
	}
}
