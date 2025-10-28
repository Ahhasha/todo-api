package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"todo-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type ListHandler struct {
	svc service.ListService
}

func NewListHandler(s service.ListService) *ListHandler {
	return &ListHandler{svc: s}
}

func (h *ListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_FAILED","message":"invalid json","details":{}}`, http.StatusBadRequest)
		return
	}

	list, err := h.svc.CreateList(ctx, req.Title, req.Description)
	if err != nil {
		http.Error(w, `{"code":"VALIDATION_FAILED","message":"title must be 1..100 chars","details":{}}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(list)
}

func (h *ListHandler) GetAllLists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l >= 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	lists, total := h.svc.GetAllListsWithPagination(ctx, limit, offset)
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(lists)
}

func (h *ListHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	list, err := h.svc.GetByID(ctx, id)
	if err != nil {
		http.Error(w, `{"code":"NOT_FOUND","message":"list not found","details":{}}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(list)
}

func (h *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_FAILED","message":"invalid json","details":{}}`, http.StatusBadRequest)
		return
	}
	list, err := h.svc.UpdateList(ctx, id, req.Title, req.Description)
	if err != nil {
		http.Error(w, `{"code":"VALIDATION_FAILED","message":"title must be 1..100 chars","details":{}}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(list)
}

func (h *ListHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.svc.Delete(ctx, id)
	if err != nil {
		http.Error(w, `{"code":"NOT_FOUND","message":"list not found","details":{}}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ListHandler) SearchLists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("title")
	if query == "" {
		http.Error(w, "missing title parameter", http.StatusBadRequest)
		return
	}

	lists, err := h.svc.SearchByTitle(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(lists) == 0 {
		http.Error(w, "list not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lists)
}
