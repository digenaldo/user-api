package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"user-api/internal/domain"
	"user-api/internal/usecase"
)

// ============================================
// HANDLER HTTP
// ============================================
// UserHandler gerencia as requisições HTTP relacionadas a usuários
// Traduz requisições HTTP para chamadas aos usecases e formata as respostas
//
// RESPONSABILIDADES DO HANDLER:
// 1. Receber requisições HTTP (ler JSON do body, parâmetros da URL)
// 2. Validar formato dos dados (JSON válido, campos presentes)
// 3. Chamar o usecase apropriado
// 4. Traduzir erros do usecase para status HTTP (400, 404, 500)
// 5. Formatar respostas JSON com status correto
//
// O QUE O HANDLER NÃO FAZ:
// - Não contém lógica de negócio (isso é do usecase)
// - Não acessa banco de dados diretamente (isso é do repository)
// - Não valida regras de negócio (ex: email válido - isso é do usecase)
type UserHandler struct {
	uc domain.UserUseCase  // Dependência: o usecase que contém a lógica de negócio
}

// NewUserHandler cria um novo handler recebendo o usecase como dependência
// Retorna *UserHandler (ponteiro) - padrão em Go para structs
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

// ============================================
// CREATE USER
// ============================================
// createUser trata requisições POST /api/v1/users
//
// SOBRE OS PARÂMETROS:
// - w http.ResponseWriter: usado para escrever a resposta HTTP
// - r *http.Request: contém informações da requisição (body, headers, etc.)
//   O * significa que é um ponteiro - Go passa por referência para evitar cópia
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	// Define uma struct anônima para receber os dados do JSON
	// As tags json:"name" mapeiam os campos do JSON para os campos da struct
	// Se o JSON tiver "name", vai para req.Name
	var req struct {
		Name  string `json:"name"`  // Campo Name mapeia para "name" no JSON
		Email string `json:"email"` // Campo Email mapeia para "email" no JSON
	}

	// Lê e decodifica o JSON do corpo da requisição
	//
	// SOBRE json.NewDecoder(r.Body).Decode(&req):
	// - r.Body é um io.Reader com os bytes do JSON enviado
	// - json.NewDecoder cria um decodificador que lê desse Reader
	// - .Decode(&req) converte o JSON para a struct req
	// - O & passa um ponteiro para req, permitindo que Decode preencha os campos
	//
	// Se o JSON for inválido (ex: sintaxe errada, tipo errado), retorna erro
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return  // Para a execução aqui - não continua
	}

	// Chama o usecase para criar o usuário
	// A validação do email (deve conter '@') acontece dentro do usecase
	//
	// CreateUser retorna (*domain.User, error)
	// - Se sucesso: user contém o usuário criado (com ID populado)
	// - Se erro: user é nil e err contém o erro
	user, err := h.uc.CreateUser(req.Name, req.Email)
	if err != nil {
		// Tratamento de erros: traduz erros do usecase para status HTTP
		// ErrInvalidEmail → 400 Bad Request (erro do cliente)
		if err == usecase.ErrInvalidEmail {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Outros erros (ex: banco indisponível) → 500 Internal Server Error
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Retorna 201 Created com o usuário criado em JSON
	// 201 Created é o status HTTP padrão para criação bem-sucedida
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
