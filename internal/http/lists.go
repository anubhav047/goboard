package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anubhav047/goboard/internal/services/list"
)

// ListHandler handles HTTP requests for lists
type ListHandler struct {
	service *list.Service
}

// NewListHandler creates a new ListHandler
func NewListHandler(service *list.Service) *ListHandler {
	return &ListHandler{
		service: service,
	}
}

// RegisterRoutes adds the list routes to router
func (h *ListHandler) RegisterRoutes(mux *http.ServeMux, mw *Middleware) {
	// All list routes require authentication
	mux.Handle("GET /api/boards/{boardId}/lists", mw.RequireAuth(http.HandlerFunc(h.handleGetBoardLists)))
	mux.Handle("POST /api/boards/{boardId}/lists", mw.RequireAuth(http.HandlerFunc(h.handleCreateList)))
	mux.Handle("GET /api/lists/{id}", mw.RequireAuth(http.HandlerFunc(h.handleGetList)))
	mux.Handle("PUT /api/lists/{id}", mw.RequireAuth(http.HandlerFunc(h.handleUpdateList)))
	mux.Handle("DELETE /api/lists/{id}", mw.RequireAuth(http.HandlerFunc(h.handleDeleteList)))
}

type CreateListRequest struct {
	Name     string `json:"name"`
	Position int32  `json:"position"`
}

type UpdateListRequest struct {
	Name     string `json:"name"`
	Position int32  `json:"position"`
}

// handleCreateList creates a new list in a board
func (h *ListHandler) handleCreateList(w http.ResponseWriter, r *http.Request) {
	// Parse board ID from URL
	boardIdStr := r.PathValue("boardId")
	boardId, err := strconv.ParseInt(boardIdStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	// Parse request body
	var req CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create list
	list, err := h.service.CreateList(r.Context(), req.Name, int32(boardId), req.Position)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, list)
}

// handleGetBoardLists gets all lists for a board
func (h *ListHandler) handleGetBoardLists(w http.ResponseWriter, r *http.Request) {
	// Parse board ID from URL
	boardIdStr := r.PathValue("boardId")
	boardId, err := strconv.ParseInt(boardIdStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	// Get board's lists
	lists, err := h.service.GetBoardLists(r.Context(), int32(boardId))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, lists)
}

// handleGetList gets a single list by ID
func (h *ListHandler) handleGetList(w http.ResponseWriter, r *http.Request) {
	// Parse list ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid list ID")
		return
	}

	// Get list
	list, err := h.service.GetListByID(r.Context(), int32(id))
	if err != nil {
		WriteError(w, http.StatusNotFound, "List not found")
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

// handleUpdateList updates a list
func (h *ListHandler) handleUpdateList(w http.ResponseWriter, r *http.Request) {
	// Parse list ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid list ID")
		return
	}

	// Parse request body
	var req UpdateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Update list
	list, err := h.service.UpdateList(r.Context(), int32(id), req.Name, req.Position)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

// handleDeleteList deletes a list
func (h *ListHandler) handleDeleteList(w http.ResponseWriter, r *http.Request) {
	// Parse list ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid list ID")
		return
	}

	// Delete list
	err = h.service.DeleteList(r.Context(), int32(id))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "List deleted successfully"})
}
