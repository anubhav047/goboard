package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anubhav047/goboard/internal/services/card"
)

// CardHandler handles HTTP requests for cards
type CardHandler struct {
	service *card.Service
}

// NewCardHandler creates a new CardHandler
func NewCardHandler(service *card.Service) *CardHandler {
	return &CardHandler{
		service: service,
	}
}

// RegisterRoutes adds the card routes to router
func (h *CardHandler) RegisterRoutes(mux *http.ServeMux, mw *Middleware) {
	// All card routes require authentication
	mux.Handle("GET /api/lists/{listId}/cards", mw.RequireAuth(http.HandlerFunc(h.handleGetListCards)))
	mux.Handle("POST /api/lists/{listId}/cards", mw.RequireAuth(http.HandlerFunc(h.handleCreateCard)))
	mux.Handle("GET /api/cards/{id}", mw.RequireAuth(http.HandlerFunc(h.handleGetCard)))
	mux.Handle("PUT /api/cards/{id}", mw.RequireAuth(http.HandlerFunc(h.handleUpdateCard)))
	mux.Handle("PUT /api/cards/{id}/move", mw.RequireAuth(http.HandlerFunc(h.handleMoveCard)))
	mux.Handle("DELETE /api/cards/{id}", mw.RequireAuth(http.HandlerFunc(h.handleDeleteCard)))
}

type CreateCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int32  `json:"position"`
}

type UpdateCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MoveCardRequest struct {
	ListID   int32 `json:"list_id"`
	Position int32 `json:"position"`
}

// handleCreateCard creates a new card in a list
func (h *CardHandler) handleCreateCard(w http.ResponseWriter, r *http.Request) {
	// Parse list ID from URL
	listIdStr := r.PathValue("listId")
	listId, err := strconv.ParseInt(listIdStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid list ID")
		return
	}

	// Parse request body
	var req CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create card
	card, err := h.service.CreateCard(r.Context(), req.Title, req.Description, int32(listId), req.Position)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, card)
}

// handleGetListCards gets all cards for a list
func (h *CardHandler) handleGetListCards(w http.ResponseWriter, r *http.Request) {
	// Parse list ID from URL
	listIdStr := r.PathValue("listId")
	listId, err := strconv.ParseInt(listIdStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid list ID")
		return
	}

	// Get list's cards
	cards, err := h.service.GetListCards(r.Context(), int32(listId))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, cards)
}

// handleGetCard gets a single card by ID
func (h *CardHandler) handleGetCard(w http.ResponseWriter, r *http.Request) {
	// Parse card ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}

	// Get card
	card, err := h.service.GetCardByID(r.Context(), int32(id))
	if err != nil {
		WriteError(w, http.StatusNotFound, "Card not found")
		return
	}

	WriteJSON(w, http.StatusOK, card)
}

// handleUpdateCard updates a card's title and description
func (h *CardHandler) handleUpdateCard(w http.ResponseWriter, r *http.Request) {
	// Parse card ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}

	// Parse request body
	var req UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Update card
	card, err := h.service.UpdateCard(r.Context(), int32(id), req.Title, req.Description)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, card)
}

// handleMoveCard moves a card to a different list and/or position (for drag & drop)
func (h *CardHandler) handleMoveCard(w http.ResponseWriter, r *http.Request) {
	// Parse card ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}

	// Parse request body
	var req MoveCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Move card
	card, err := h.service.MoveCard(r.Context(), int32(id), req.ListID, req.Position)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, card)
}

// handleDeleteCard deletes a card
func (h *CardHandler) handleDeleteCard(w http.ResponseWriter, r *http.Request) {
	// Parse card ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}

	// Delete card
	err = h.service.DeleteCard(r.Context(), int32(id))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Card deleted successfully"})
}
