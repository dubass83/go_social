package main

import (
	"github.com/dubass83/go_social/internal/env"
	"github.com/dubass83/go_social/internal/store"
	"github.com/rs/zerolog/log"
)

func main() {
	conf := config{
		addr: env.GetString("API_ADDR", ":8080"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: conf,
		store:  store,
	}

	if err := app.run(app.mount()); err != nil {
		log.Fatal().Err(err).Msg("failed to run application")
	}
}
