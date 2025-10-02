package main

import (
	"os"

	"github.com/dubass83/go_social/internal/db"
	"github.com/dubass83/go_social/internal/env"
	"github.com/dubass83/go_social/internal/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const version = "0.1.0"

//	@title			GO Social Study App
//	@description	This is a sample server Go Social server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	conf := config{
		addr:   env.GetString("API_ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		env:    env.GetString("ENVIRONMENT", "dev"),
		db: dbConf{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:password@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10m"),
		},
	}

	// Logger
	if conf.env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// log.Debug().Msgf("config values: %+v", conf)
	}

	db, err := db.New(
		conf.db.addr,
		conf.db.maxOpenConns,
		conf.db.maxIdleConns,
		conf.db.maxIdleTime,
	)

	if err != nil {
		log.Debug().Msgf("db config: %v", conf.db)
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
