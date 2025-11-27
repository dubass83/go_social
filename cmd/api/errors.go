package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().Err(err).Msgf("internal server error: %s path: %s", r.Method, r.RequestURI)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered an error")
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().Err(err).Msgf("bad request error: %s path: %s", r.Method, r.RequestURI)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().Err(err).Msgf("not found error: %s path: %s", r.Method, r.RequestURI)
	writeJSONError(w, http.StatusNotFound, "not found")
}

func unAuthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().Err(err).Msgf("unauthorized error: %s path: %s", r.Method, r.RequestURI)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
