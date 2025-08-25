package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
)

type PostPayload struct {
	Title   string   `json:"title" validate:"required,min=2,max=100"`
	Content string   `json:"content" validate:"required,min=2,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload PostPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err = validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: add user ID from auth
		UserID: 1,
	}

	if err := app.store.Post.Create(r.Context(), post); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (app *application) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	ctx := r.Context()

	post, err := app.store.Post.GetByID(ctx, postID)
	if err != nil {
		if err == store.ErrNotFound {
			notFoundResponse(w, r, err)
			return
		}
		internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}
