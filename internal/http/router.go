package http

import (
	"net/http"
	"todo-api/internal/http/handlers"
	"todo-api/internal/http/middleware"
	"todo-api/internal/service"
	"todo-api/internal/storage/mem"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	listRepo := mem.NewListRepo()
	listService := service.NewListService(listRepo)
	listHandler := handlers.NewListHandler(listService)

	r.Use(middleware.Logging)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/lists", func(r chi.Router) {
			r.Post("/", listHandler.CreateList)
			r.Get("/", listHandler.GetAllLists)
			r.Get("/{id}", listHandler.GetByID)
			r.Patch("/{id}", listHandler.UpdateTitle)
			r.Delete("/{id}", listHandler.Delete)
		})

	})

	r.Get("/openapi.yaml", handlers.OpenAPISpec)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return r
}
