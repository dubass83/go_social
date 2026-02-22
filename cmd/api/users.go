package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
)

type UserKey string

const UserCtxKey UserKey = "user"

// GetUserByIDHandler godoc
//
//	@Summary		Get a user
//	@Description	get user by ID
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id} [get]
func (app *application) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		internalServerError(w, r, err)
	}

}

// GetUsersPostsHandler godoc
//
//	@Summary		Get a users posts
//	@Description	Fetch all posts authored by a specific user, for display on their profile page.
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Param			limit	query		int		false	"Limit number of posts"	default(10)
//	@Param			offset	query		int		false	"Offset for pagination"	default(0)
//	@Param			sort	query		string	false	"Sort order (asc/desc)"	default(desc)
//	@Param			tags	query		string	false	"Filter by tags (comma-separated)"
//	@Param			search	query		string	false	"Search query"
//	@Success		200	{object}	[]*store.PostWithMetadata
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id}/posts [get]
func (app *application) GetUsersPostsHandler(w http.ResponseWriter, r *http.Request) {
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

	user := getUserFromCtx(r)

	posts, err := app.store.Post.GetUserPosts(r.Context(), user.ID, pgPostsQuery)
	if err != nil {
		internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		internalServerError(w, r, err)
	}

}

// activateUserHandler godoc
//
//	@Summary		Activate a user
//	@Description	activate user by invitation token
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation Token"
//	@Success		202		{string}	string	"User activated successfully"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put,get]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.User.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			badRequestResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}

	err = app.jsonResponse(w, http.StatusAccepted, "User activated successfully")
	if err != nil {
		internalServerError(w, r, err)
	}
}

// func (app *application) userContextMiddelware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
// 		if err != nil {
// 			badRequestResponse(w, r, err)
// 			return
// 		}

// 		user, err := app.store.User.GetByID(r.Context(), userID)
// 		if err != nil {
// 			switch err {
// 			case store.ErrNotFound:
// 				notFoundResponse(w, r, err)
// 				return
// 			default:
// 				internalServerError(w, r, err)
// 				return
// 			}
// 		}

// 		ctx := context.WithValue(r.Context(), UserCtxKey, user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func getUserFromCtx(r *http.Request) *store.User {
	user := r.Context().Value(UserCtxKey).(*store.User)
	return user
}
