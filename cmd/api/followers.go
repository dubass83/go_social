package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FollowerIDPayload struct {
	ID int64 `json:"id"`
}

// FollowUserByIDHandler godoc
//
//	@Summary		Follow a user
//	@Description	follow user by ID
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		202
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{userID}/follow [put]
func (app *application) FollowUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	flID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	user := getUserFromCtx(r)

	if err := app.store.Follower.CreateFollower(r.Context(), flID, user.ID); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusAccepted, nil); err != nil {
		internalServerError(w, r, err)
	}
}

// UnfollowUserByIDHandler godoc
//
//	@Summary		Unfollow a user
//	@Description	unfollow user by ID
//	@Tags			USERS
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		202
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{userID}/unfollow [put]
func (app *application) UnfollowUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	flID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	user := getUserFromCtx(r)

	if err := app.store.Follower.DeleteFollower(r.Context(), flID, user.ID); err != nil {
		internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusAccepted, nil); err != nil {
		internalServerError(w, r, err)
	}
}
