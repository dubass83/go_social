package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}
	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
