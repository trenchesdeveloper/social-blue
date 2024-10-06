package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/trenchesdeveloper/social-blue/config"
	"go.uber.org/zap"

	"github.com/trenchesdeveloper/social-blue/docs" //This is required for swaggo to find your docs
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"net/http"
	"time"
)

type server struct {
	config *config.AppConfig
	store  db.Store
	logger *zap.SugaredLogger
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
		docsURL := fmt.Sprintf("%s/swagger/doc.json", s.config.ServerPort)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL), //The url pointing to API definition
		))

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
				r.Use(s.userContextMiddleware)
				r.Get("/", s.getUserHandler)
				r.Put("/follow", s.followUserHandler)
				r.Put("/unfollow", s.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", s.GetUserFeedsHandler)
			})
		})
	})
	return r
}

func (s *server) start(mux http.Handler) error {
	//Docs

	docs.SwaggerInfo.Title = "Social Blue API"
	docs.SwaggerInfo.Description = "This is a social media API"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = s.config.ApiUrl

	srv := &http.Server{
		Addr:         s.config.ServerPort,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	s.logger.Infow("Starting server", "port", s.config.ServerPort, "env", s.config.Environment)
	return srv.ListenAndServe()
}
