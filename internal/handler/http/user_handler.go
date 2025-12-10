package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"user-api/internal/domain"
	"user-api/internal/usecase"
)

// UserHandler handles HTTP requests for users
// Comentários gerais:
//   - Este handler expõe os endpoints REST para CRUD de usuários.
//   - Cada método traduz a requisição HTTP para chamadas aos usecases (lógica de negócio),
//     trata erros e escreve respostas JSON com o status HTTP apropriado.
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

	// Decodifica o JSON do corpo da requisição para a struct local `req`.
	// Comentário linha-a-linha (trecho importante):
	// - `r.Body` é um io.Reader com os bytes enviados pelo cliente
	// - `json.NewDecoder(...).Decode(&req)` faz o parse do JSON para a struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Se o JSON for inválido retornamos 400 Bad Request.
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Chama o usecase para criar o usuário. A lógica de validação de email
	// está encapsulada no usecase (ver internal/usecase/user_usecase.go).
	user, err := h.uc.CreateUser(req.Name, req.Email)
	if err != nil {
		// Tratamento de erros de negócio: se o email for inválido retornamos 400.
		if err == usecase.ErrInvalidEmail {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Qualquer outro erro é considerado erro interno do servidor.
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Sucesso: retorna 201 Created com o usuário criado em JSON.
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

	// `chi.URLParam` extrai o parâmetro `id` da rota.
	// O usecase espera um ID no formato hexadecimal de ObjectID do MongoDB.
	user, err := h.uc.GetUser(id)
	if err != nil {
		// Se o usecase sinalizar que o usuário não foi encontrado, retornamos 404.
		if err == usecase.ErrNotFound {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		// Erro genérico do servidor
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

	// Decodifica o corpo JSON para os campos opcionais de atualização.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Chama o usecase para atualizar o usuário. O usecase lida com:
	// - verificação da existência do usuário
	// - validação do e-mail (se informado)
	// - persistência via repositório
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
