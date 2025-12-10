package http

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "user-api/docs" // Importa o pacote docs gerado pelo swag init
)

// RegisterSwagger registra as rotas do Swagger UI
// A documentação interativa estará disponível em /swagger/index.html
func RegisterSwagger(r chi.Router) {
	// WrapHandler serve automaticamente os arquivos do pacote docs
	// quando ele está importado (linha acima)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}

