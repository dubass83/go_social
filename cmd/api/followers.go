package main

import "net/http"

type FollowerIDPayload struct {
	ID int64 `json:"id"`
}

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
