package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"user-api/internal/domain"
	"user-api/internal/usecase"
)

// UserHandler gerencia as requisições HTTP relacionadas a usuários
// Traduz requisições HTTP para chamadas aos usecases e formata as respostas
type UserHandler struct {
	uc domain.UserUseCase
}

// NewUserHandler cria um novo handler recebendo o usecase como dependência
func NewUserHandler(uc domain.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// RegisterRoutes registra todas as rotas de usuários no router
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/", h.listUsers)
		r.Get("/{id}", h.getUser)
		r.Put("/{id}", h.updateUser)
		r.Delete("/{id}", h.deleteUser)
	})
}

// createUser trata requisições POST /api/v1/users
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Lê o JSON do corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Chama o usecase para criar o usuário
	// A validação do email acontece dentro do usecase
	user, err := h.uc.CreateUser(req.Name, req.Email)
	if err != nil {
		if err == usecase.ErrInvalidEmail {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Retorna 201 Created com o usuário criado
	writeJSON(w, http.StatusCreated, user)
}

// listUsers trata requisições GET /api/v1/users
func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.ListUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// getUser trata requisições GET /api/v1/users/{id}
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

// updateUser trata requisições PUT /api/v1/users/{id}
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

// deleteUser trata requisições DELETE /api/v1/users/{id}
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

	// DELETE bem-sucedido retorna 204 No Content (sem corpo)
	w.WriteHeader(http.StatusNoContent)
}

// writeJSON escreve uma resposta JSON com o status HTTP informado
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError escreve uma resposta de erro em JSON
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
