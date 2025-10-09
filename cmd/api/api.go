package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dubass83/go_social/docs"
	"github.com/dubass83/go_social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	addr   string
	apiURL string
	env    string
	db     dbConf
}

type dbConf struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type registerUserPayload struct {
	Username string `json:"username" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
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
		docsURL := fmt.Sprintf("http://%s/v1/swagger/doc.json", app.config.apiURL)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

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
			r.Put("/activate/{token}", app.activateUserHandler)
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

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

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
