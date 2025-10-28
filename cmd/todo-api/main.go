package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo-api/internal/config"
	"todo-api/internal/database"
	httphandlers "todo-api/internal/http"
	"todo-api/internal/http/handlers"
	"todo-api/internal/service"
	"todo-api/internal/storage/postgres"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	log.Println("Connecting to database...")
	pool, err := database.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to database")

	repo := postgres.NewListRepo(pool)
	taskRepo := postgres.NewTaskRepo(pool)

	svc := service.NewListService(repo)
	taskSvc := service.NewTaskService(taskRepo, repo)

	listHandler := handlers.NewListHandler(svc)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	router := httphandlers.NewRouter(listHandler, taskHandler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
