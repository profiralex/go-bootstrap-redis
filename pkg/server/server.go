package server

import (
	"fmt"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/profiralex/go-bootstrap-redis/pkg/config"
	"github.com/profiralex/go-bootstrap-redis/pkg/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve(cfg config.Config) error {
	render.Respond = NewApiRenderer()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Group(func(r chi.Router) {
		r.Use(db.SessionCtx)

		entitiesController := newEntitiesController()
		r.Route("/v1", func(r chi.Router) {
			r.Post("/entities", entitiesController.createEntity)
			r.Get("/entities/{uuid}", entitiesController.getEntity)
		})
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppConfig.Port), r)
	if err != nil {
		return fmt.Errorf("server failed %w", err)
	}

	return nil
}
