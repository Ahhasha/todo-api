package http

import (
	"net/http"

	"todo-api/internal/http/handlers"
	"todo-api/internal/http/middleware"
	"todo-api/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(listService service.ListService) http.Handler {
	r := chi.NewRouter()

	listHandler := handlers.NewListHandler(listService)

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/lists", func(r chi.Router) {
			r.Post("/", listHandler.CreateList)
			r.Get("/", listHandler.GetAllLists)
			r.Get("/search", listHandler.SearchLists)
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
