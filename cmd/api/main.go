package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	httphandler "user-api/internal/handler/http"
	"user-api/internal/infra/mongo"
	"user-api/internal/repository"
	"user-api/internal/usecase"
)

func main() {
	// Lê variáveis de ambiente ou usa valores padrão
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Conecta ao MongoDB
	// O NewClient retorna um ponteiro (*mongo.Client) - isso significa que
	// todas as operações usam o mesmo cliente compartilhado, sem cópias
	client := mongo.NewClient(mongoURI)
	
	// Garante que desconecta do MongoDB quando a aplicação encerrar
	defer func() {
		if err := client.Disconnect(nil); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Seleciona o database que vamos usar
	db := client.Database("userdb")

	// Monta a cadeia de dependências: repository → usecase → handler
	// Cada camada recebe a anterior como dependência (injeção de dependência)
	repo := repository.NewUserMongoRepository(db)
	uc := usecase.NewUserUseCase(repo)
	handler := httphandler.NewUserHandler(uc)

	// Configura as rotas HTTP
	r := chi.NewRouter()
	httphandler.RegisterHealth(r)  // Rota de healthcheck
	handler.RegisterRoutes(r)      // Rotas de usuários (CRUD)

	// Inicia o servidor HTTP
	// Em produção, considere usar um servidor com timeouts e TLS configurados
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
