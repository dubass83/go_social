package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	conf := config{
		addr: ":8080",
	}
	app := &application{
		config: conf,
	}

	if err := app.run(app.mount()); err != nil {
		log.Fatal().Err(err).Msg("failed to run application")
	}
}
