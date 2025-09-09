package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
)

type commentPayLoad struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

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
