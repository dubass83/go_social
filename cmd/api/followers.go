package main

import "net/http"

type FollowerIDPayload struct {
	ID int64 `json:"id"`
}

// FollowUserByIDHandler godoc
//
//	@Summary		Follow a user
//	@Description	follow user by ID
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Posts ID"
//	@Success		202
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id}/follow [put]
func (app *application) FollowUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	fl := FollowerIDPayload{}
	if err := readJSON(w, r, &fl); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user := getUserFromCtx(r)

	if err := app.store.Follower.CreateFollower(r.Context(), user.ID, fl.ID); err != nil {
		internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusAccepted, nil); err != nil {
		internalServerError(w, r, err)
	}
}

// UnfollowUserByIDHandler godoc
//
//	@Summary		Unfollow a user
//	@Description	unfollow user by ID
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Posts ID"
//	@Success		202
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id}/unfollow [put]
func (app *application) UnfollowUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	fl := FollowerIDPayload{}
	if err := readJSON(w, r, &fl); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user := getUserFromCtx(r)

	if err := app.store.Follower.DeleteFollower(r.Context(), user.ID, fl.ID); err != nil {
		internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusAccepted, nil); err != nil {
		internalServerError(w, r, err)
	}
}
