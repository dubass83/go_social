package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// HealthCheckHandler godoc
//
//	@Summary		Health Check
//	@Description	get health status
//	@Tags			OPS
//	@Produce		json
//	@Success		200	{map}	map[string]string
//	@Failure		500	{map}	map[string]string
//	@Router			/health [get]
func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}
	err := app.jsonResponse(w, http.StatusOK, data)
	if err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
