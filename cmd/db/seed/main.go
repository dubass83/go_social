package main

import (
	"github.com/dubass83/go_social/internal/db"
	"github.com/dubass83/go_social/internal/env"
	"github.com/dubass83/go_social/internal/store"
	"github.com/rs/zerolog/log"
)

type config struct {
	db dbConf
}

type dbConf struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func main() {
	conf := config{
		db: dbConf{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:password@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10m"),
		},
	}

	conn, err := db.New(
		conf.db.addr,
		conf.db.maxOpenConns,
		conf.db.maxIdleConns,
		conf.db.maxIdleTime,
	)

	if err != nil {
		log.Debug().Msgf("db config: %v", conf.db)
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer conn.Close()
	log.Info().Msg("database connection established")

	store := store.NewStorage(conn)

	db.Seed(store, 100)

	log.Info().Msg("database seeded")
}
