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

// CreatePostHandler godoc
//
//	@Summary		Create a new post
//	@Description	create a new post with title, content and tags
//	@Tags			POSTS
//	@Accept			json
//	@Produce		json
//	@Param			post	body		PostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/posts [post]
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

	user := getUserFromCtx(r)

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
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

// GetAllPostsHandler godoc
//
//	@Summary		Get a paginated list of all posts
//	@Description	Return a paginated list of all posts (not filtered by follows), useful as a public discovery feed for new users who do not follow anyone yet.
//	@Tags			POSTS
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit number of posts"	default(10)
//	@Param			offset	query		int		false	"Offset for pagination"	default(0)
//	@Param			sort	query		string	false	"Sort order (asc/desc)"	default(desc)
//	@Param			tags	query		string	false	"Filter by tags (comma-separated)"
//	@Param			search	query		string	false	"Search query"
//	@Success		200		{object}    []*store.PostWithMetadata
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/posts [get]
func (app *application) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	pgPostsQueryDefault := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	pgPostsQuery, err := pgPostsQueryDefault.Parse(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(pgPostsQuery); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	posts, err := app.store.Post.GetAllPosts(ctx, pgPostsQuery)

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		internalServerError(w, r, err)
		return
	}
}

// GetPostByIDHandler godoc
//
//	@Summary		Show a post with comments
//	@Description	get post by ID
//	@Tags			POSTS
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Posts ID"
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/posts/{id} [get]
func (app *application) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {

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

// DeletePostHandler godoc
//
//	@Summary		Delete a post
//	@Description	delete post by ID
//	@Tags			POSTS
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/posts/{id} [delete]
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

// UpdatePostHandler godoc
//
//	@Summary		Update a post
//	@Description	update post by ID with optional title, content and tags
//	@Tags			POSTS
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			post	body		UpdatePostPayload	true	"Update post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/posts/{id} [patch]
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

	user := getUserFromCtx(r)

	// if user.ID != post.UserID {
	// 	unAuthorizedResponse(w, r, fmt.Errorf("user with ID %d is not allowed to update this post", user.ID))
	// 	return
	// }

	updatedPost := &store.Post{
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
		UserID:  user.ID,
	}

	if err := app.store.Post.Update(ctx, post.ID, post.Version, updatedPost); err != nil {
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
