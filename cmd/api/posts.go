package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCTX postKey = "post"

type PostPayload struct {
	Title   string   `json:"title" validate:"required,min=2,max=100"`
	Content string   `json:"content" validate:"required,min=2,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   string   `json:"title" validate:"omitempty,min=2,max=100"`
	Content string   `json:"content" validate:"omitempty,min=2,max=1000"`
	Tags    []string `json:"tags" validate:"omitempty"`
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

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (app *application) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	post := getPostFromCtx(r)

	comments, err := app.store.Comment.GetByPostID(ctx, post.ID)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (app *application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	ctx := r.Context()

	if err := app.store.Post.DeleteByID(ctx, postID); err != nil {
		if err == store.ErrNotFound {
			notFoundResponse(w, r, err)
			return
		}
		internalServerError(w, r, err)
		return
	}
	data := map[string]string{
		"message": fmt.Sprintf("post with id %s was successfully deleted from the database", postID),
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if payload.Content != "" {
		post.Content = payload.Content
	}
	if payload.Title != "" {
		post.Title = payload.Title
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	updatedPost := &store.Post{
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
		// TODO: add user ID from auth
		UserID: 1,
	}

	if err := app.store.Post.Update(ctx, post.ID, updatedPost); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedPost); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (app *application) postContextMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		ctx = context.WithValue(ctx, postCTX, post)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCTX).(*store.Post)
	return post
}
