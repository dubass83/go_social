package main

import (
	"github.com/dubass83/go_social/internal/db"
	"github.com/dubass83/go_social/internal/env"
	"github.com/dubass83/go_social/internal/store"
	"github.com/rs/zerolog/log"
)

func main() {
	conf := config{
		addr: env.GetString("API_ADDR", ":8080"),
		db: dbConf{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:password@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10m"),
		},
	}

	db, err := db.New(
		conf.db.addr,
		conf.db.maxOpenConns,
		conf.db.maxIdleConns,
		conf.db.maxIdleTime,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()
	log.Info().Msg("database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: conf,
		store:  store,
	}

	if err := app.run(app.mount()); err != nil {
		log.Fatal().Err(err).Msg("failed to run application")
	}
}
