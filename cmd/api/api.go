package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/trenchesdeveloper/social-blue/config"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type server struct {
	config *config.AppConfig
	store  db.Store
}

func (s *server) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", s.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", s.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(s.postsContextMiddleware)
				r.Get("/", s.getPostHandler)
				r.Patch("/", s.updatePostHandler)
				r.Delete("/", s.deletePostHandler)
			})

			//r.Put("/{id}", s.updatePost)
			//r.Delete("/{id}", s.deletePost)
		})

		r.Route("/users", func(r chi.Router) {
			//r.Post("/", s.createUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", s.getUserHandler)
			})
		})
	})

	return r
}

func (s *server) start(mux http.Handler) error {
	srv := &http.Server{
		Addr:         s.config.ServerPort,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server starting on %s", s.config.ServerPort)
	return srv.ListenAndServe()
}
