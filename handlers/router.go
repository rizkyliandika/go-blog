package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Response struct {
	Msg  string
	Code int
	Data any
}

func CreateRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api", func(r chi.Router) {
		// version 1
		r.Route("/v1", func(r chi.Router) {
			r.Get("/healthCheck", healthCheck)
			r.Get("/post", getPosts)
			r.Get("/post/{id}", getPostById)
			r.Post("/post/create", createPost)
			r.Put("/post/update", updatePost)
			r.Delete("/post/delete/{id}", deletePost)
		})
	})

	return router
}
