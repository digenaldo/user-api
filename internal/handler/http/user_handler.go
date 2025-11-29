package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"user-api/internal/domain"
	"user-api/internal/usecase"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	uc domain.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(uc domain.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/", h.listUsers)
		r.Get("/{id}", h.getUser)
		r.Put("/{id}", h.updateUser)
		r.Delete("/{id}", h.deleteUser)
	})
}

// createUser handles POST /api/v1/users
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.uc.CreateUser(req.Name, req.Email)
	if err != nil {
		if err == usecase.ErrInvalidEmail {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// listUsers handles GET /api/v1/users
func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.ListUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// getUser handles GET /api/v1/users/{id}
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.uc.GetUser(id)
	if err != nil {
		if err == usecase.ErrNotFound {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// updateUser handles PUT /api/v1/users/{id}
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.uc.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		if err == usecase.ErrNotFound {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		if err == usecase.ErrInvalidEmail {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// deleteUser handles DELETE /api/v1/users/{id}
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.uc.DeleteUser(id)
	if err != nil {
		if err == usecase.ErrNotFound {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
