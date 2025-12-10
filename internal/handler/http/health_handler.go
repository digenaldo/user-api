package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// RegisterHealth registra a rota de healthcheck usada por sistemas de
// monitoramento e para testes manuais simples.
// Rota: GET /healthz
func RegisterHealth(r chi.Router) {
	r.Get("/healthz", healthz)
}

// healthz é o handler que responde com um JSON simples contendo o status
// e um carimbo de tempo UTC. Este endpoint deve ser rápido e não deve
// executar consultas pesadas; seu objetivo é indicar se a aplicação está
// inicializada e capaz de responder a requisições HTTP.
func healthz(w http.ResponseWriter, r *http.Request) {
	// Define o content-type JSON
	w.Header().Set("Content-Type", "application/json")
	// Retorna HTTP 200 OK
	w.WriteHeader(http.StatusOK)
	// Encode do payload JSON. Em produção você pode incluir checagens
	// adicionais (por exemplo, ping ao banco) — mas cuidado para não
	// transformar este endpoint em uma operação lenta.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
