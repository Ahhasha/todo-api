package http

import (
	"net/http"

	_ "todo-api/docs"
	"todo-api/internal/http/handlers"
	"todo-api/internal/http/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(listHandler *handlers.ListHandler, taskHandler *handlers.TaskHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/lists", func(r chi.Router) {
			r.Post("/", listHandler.CreateList)
			r.Get("/", listHandler.GetAllLists)
			r.Get("/search", listHandler.SearchLists)
			r.Get("/{id}", listHandler.GetByID)
			r.Patch("/{id}", listHandler.UpdateList)
			r.Delete("/{id}", listHandler.Delete)
		})
		r.Route("/lists/{listID}/tasks", func(r chi.Router) {
			r.Post("/", taskHandler.CreateTask)
			r.Get("/", taskHandler.ListTasks)
		})

	})

	r.Route("/api/v1/tasks/{taskID}", func(r chi.Router) {
		r.Get("/", taskHandler.GetTask)
		r.Patch("/", taskHandler.UpdateTask)
		r.Delete("/", taskHandler.DeleteTask)
	})

	r.Get("/openapi.yaml", handlers.OpenAPISpec)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/openapi.yaml"),
	))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return r
}
