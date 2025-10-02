package main

import (
	"log"
	"net/http"
	"os"
	apphttp "todo-api/internal/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := apphttp.NewRouter()

	log.Printf("Server is running on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
