package handlers

import (
	"net/http"
)

func OpenAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/openapi.yaml")
}
