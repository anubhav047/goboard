package http

import (
	"encoding/json"
	"net/http"

	"github.com/anubhav047/goboard/internal/services/user"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	service *user.Service
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// RegisterRoutes adds the user routes to router.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.handleRegister)
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handleRegister is the handler for the user registration endpoint.
func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Execute Business Logic and passing request context all the way to service
	user, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		// TODO : Add specific checks like duplicate email
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Encode and send response
	response := map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}
	WriteJSON(w, http.StatusCreated, response)
}
