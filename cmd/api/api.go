package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) run() error {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}
	log.Info().Msgf("Starting server on %s", app.config.addr)
	return srv.ListenAndServe()
}
