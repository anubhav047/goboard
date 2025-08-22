package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/anubhav047/goboard/internal/services/user"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	service *user.Service
	sm      *scs.SessionManager
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *user.Service, sm *scs.SessionManager) *UserHandler {
	return &UserHandler{
		service: service,
		sm:      sm,
	}
}

// RegisterRoutes adds the user routes to router.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.handleRegister)
	mux.Handle("POST /login", h.sm.LoadAndSave(http.HandlerFunc(h.handleLogin)))
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	// Decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Execute the login logic.
	userr, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// Check if it's our specific invalid credentials error.
		if errors.Is(err, user.ErrInvalidCredentials) {
			WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		// For all other errors, return a generic server error.
		WriteError(w, http.StatusInternalServerError, "An unexpected error occurred")
		return
	}

	// Use the session manager from the handler struct.
	err = h.sm.RenewToken(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to renew session token")
		return
	}
	h.sm.Put(r.Context(), "authenticatedUserID", userr.ID)

	response := map[string]interface{}{
		"id":    userr.ID,
		"name":  userr.Name,
		"email": userr.Email,
	}
	WriteJSON(w, http.StatusOK, response)
}
