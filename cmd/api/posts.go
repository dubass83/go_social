package main

import (
	"net/http"

	"github.com/dubass83/go_social/internal/store"
	"github.com/rs/zerolog/log"
)

type PostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload PostPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		log.Error().Err(err).Msg("failed to read JSON")
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: add user ID from auth
		UserID: 1,
	}

	if err := app.store.Post.Create(r.Context(), post); err != nil {
		log.Error().Err(err).Msg("failed to create post")
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		log.Error().Err(err).Msg("failed to write JSON")
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
