package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// RegisterHealth registra a rota de healthcheck
// Útil para monitoramento e verificar se a aplicação está respondendo
func RegisterHealth(r chi.Router) {
	r.Get("/healthz", healthz)
}

// healthz retorna um JSON simples indicando que a aplicação está funcionando
// Este endpoint deve ser rápido - não faça consultas pesadas aqui
// Em produção, você pode adicionar checagens (ex: ping ao banco), mas mantenha rápido
//
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthz [get]
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
