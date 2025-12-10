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
	// OBS: Este arquivo é o entrypoint da aplicação.
	// Comentários abaixo explicam cada etapa de inicialização do servidor HTTP,
	// a conexão com o MongoDB e o registro das rotas.

	// Read environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create MongoDB client
	// Cria o cliente MongoDB usando o helper em internal/infra/mongo.
	// O helper encapsula opções do driver e retorna um `*mongo.Client` pronto.
	//
	// Explicando `*` no contexto de `*mongo.Client`:
	// - Em Go, `*T` significa "ponteiro para T". Quando a função retorna
	//   `*mongo.Client` ela está retornando o endereço de um valor do tipo
	//   `mongo.Client` alocado (ou gerenciado) internamente.
	// - Usamos ponteiros para evitar cópias da estrutura e para que chamadas
	//   posteriores (`client.Database(...)`, `client.Disconnect(...)`) atuem
	//   sobre o mesmo cliente compartilhado.
	client := mongo.NewClient(mongoURI)
	defer func() {
		if err := client.Disconnect(nil); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Get database
	// Seleciona o database utilizado pela aplicação (userdb).
	// Todas as operações de repositório usam este database.
	db := client.Database("userdb")

	// Create repository
	// Cria a implementação do repositório que conversa com o MongoDB.
	// A camada de repositório traduz chamadas do caso de uso para operações
	// na collection `users` do MongoDB.
	repo := repository.NewUserMongoRepository(db)

	// Create use case
	// Instancia a camada de uso (usecase) passando o repositório.
	// Os usecases contêm a lógica de negócio e orquestram validações e
	// chamadas ao repositório.
	uc := usecase.NewUserUseCase(repo)

	// Create handler
	// Cria os handlers HTTP que irão expor os endpoints da aplicação.
	// Os handlers recebem a camada de usecase para delegar operações.
	handler := httphandler.NewUserHandler(uc)

	// Setup router
	r := chi.NewRouter()
	// Registra rota de healthcheck (ponto simples para monitoramento)
	httphandler.RegisterHealth(r)
	// Registra as rotas de usuário (CRUD) em /api/v1/users
	handler.RegisterRoutes(r)

	// Start server
	// Inicia o servidor HTTP. Em produção você pode substituir por um servidor
	// mais robusto (com timeouts e TLS). Aqui usamos ListenAndServe direto
	// para simplicidade didática.
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
