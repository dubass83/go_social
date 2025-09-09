package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
)

type commentPayLoad struct {
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}

func (app *application) CreateCommentToPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getPostFromCtx(r)

	payload := commentPayLoad{}

	if err := readJSON(w, r, &payload); err != nil {
		internalServerError(w, r, err)
	}

	comment := &store.Comment{
		UserID:  payload.UserID,
		Content: payload.Content,
		PostID:  post.ID,
	}
	err := app.store.Comment.Create(ctx, comment)
	if err != nil {
		internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		internalServerError(w, r, err)
		return
	}
}
