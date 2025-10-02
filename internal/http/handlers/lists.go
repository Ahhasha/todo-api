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
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_FAILED","message":"invalid json","details":{}}`, http.StatusBadRequest)
		return
	}

	list, err := h.svc.CreateList(ctx, req.Title)
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

func (h *ListHandler) UpdateTitle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"VALIDATION_FAILED","message":"invalid json","details":{}}`, http.StatusBadRequest)
		return
	}
	list, err := h.svc.UpdateTitle(ctx, id, req.Title)
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
