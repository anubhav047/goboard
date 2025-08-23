package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anubhav047/goboard/internal/db"
	"github.com/anubhav047/goboard/internal/services/board"
)

// BoardHandler handles HTTP requests for boards
type BoardHandler struct {
	service *board.Service
}

// NewBoardHandler creates a new BoardHandler
func NewBoardHandler(service *board.Service) *BoardHandler {
	return &BoardHandler{
		service: service,
	}
}

// RegisterRoutes adds the board routes to router
func (h *BoardHandler) RegisterRoutes(mux *http.ServeMux, mw *Middleware) {
	// All Board routes require authentication
	mux.Handle("GET /api/boards", mw.RequireAuth((http.HandlerFunc(h.handleGetUserBoards))))
	mux.Handle("POST /api/boards", mw.RequireAuth(http.HandlerFunc(h.handleCreateBoard)))
	mux.Handle("GET /api/boards/{id}", mw.RequireAuth(http.HandlerFunc(h.handleGetBoard)))
	mux.Handle("PUT /api/boards/{id}", mw.RequireAuth(http.HandlerFunc(h.handleUpdateBoard)))
	mux.Handle("DELETE /api/boards/{id}", mw.RequireAuth(http.HandlerFunc(h.handleDeleteBoard)))
}

type CreateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// handleCreateBoard creates a new board
func (h *BoardHandler) handleCreateBoard(w http.ResponseWriter, r *http.Request) {
	// Get User from context (set by RequireAuth middleware)
	user, ok := r.Context().Value(userContextKey).(db.User)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "Error retrieving user from context")
		return
	}

	// Parse request body
	var req CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	// Create Board
	board, err := h.service.CreateBoard(r.Context(), req.Name, req.Description, user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return created board
	WriteJSON(w, http.StatusCreated, board)
}

// handleGetUserBoards gets all boards for the authenticated user
func (h *BoardHandler) handleGetUserBoards(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := r.Context().Value(userContextKey).(db.User)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "Error retrieving user from context")
		return
	}

	// Get user's boards
	boards, err := h.service.GetUserBoards(r.Context(), user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, boards)
}

// handleGetBoard gets a single board by ID
func (h *BoardHandler) handleGetBoard(w http.ResponseWriter, r *http.Request) {
	// Parse board ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	// Get board
	board, err := h.service.GetBoardByID(r.Context(), int32(id))
	if err != nil {
		WriteError(w, http.StatusNotFound, "Board not found")
		return
	}

	WriteJSON(w, http.StatusOK, board)
}

// handleUpdateBoard updates a board
func (h *BoardHandler) handleUpdateBoard(w http.ResponseWriter, r *http.Request) {
	// Parse board ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	// Parse request body
	var req UpdateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Update board
	board, err := h.service.UpdateBoard(r.Context(), int32(id), req.Name, req.Description)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, board)
}

// handleDeleteBoard deletes a board
func (h *BoardHandler) handleDeleteBoard(w http.ResponseWriter, r *http.Request) {
	// Parse board ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	// Delete board
	err = h.service.DeleteBoard(r.Context(), int32(id))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Board deleted successfully"})
}
