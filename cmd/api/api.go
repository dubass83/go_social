package main

import (
	"net/http"
	"time"

	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	addr string
	db   dbConf
}

type dbConf struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.HealthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.CreatePostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddelware)

				r.Get("/", app.GetPostByIDHandler)
				r.Delete("/", app.DeletePostHandler)
				r.Patch("/", app.UpdatePostHandler)
				r.Post("/comments", app.CreateCommentToPostByIDHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.CreateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddelware)

				r.Get("/", app.GetUserByIDHandler)
				r.Put("/follow", app.FollowUserByIDHandler)
				r.Put("/unfollow", app.UnfollowUserByIDHandler)
				// r.Delete("/", app.DeletePostHandler)
				// r.Patch("/", app.UpdatePostHandler)
				// r.Post("/comments", app.CreateCommentToPostByIDHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.GetUserFeedHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Info().Msgf("starting server on %s", app.config.addr)
	return srv.ListenAndServe()
}
