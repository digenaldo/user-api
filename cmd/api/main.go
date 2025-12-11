// @title User API
// @version 1.0
// @description API REST de exemplo para CRUD de usuários usando Go e MongoDB
// @host localhost:8080
// @BasePath /
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
	// ============================================
	// CONFIGURAÇÃO INICIAL
	// ============================================
	// Lê variáveis de ambiente ou usa valores padrão
	// os.Getenv() retorna uma string vazia se a variável não existir
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// ============================================
	// CONEXÃO COM MONGODB
	// ============================================
	// NewClient retorna um ponteiro (*mongo.Client)
	//
	// O QUE É UM PONTEIRO?
	// - Um ponteiro é o endereço de memória onde um valor está armazenado
	// - Em Go, *T significa "ponteiro para tipo T"
	// - O operador & cria um ponteiro (pega o endereço de um valor)
	// - O operador * desreferencia um ponteiro (acessa o valor apontado)
	//
	// POR QUE USAR PONTEIROS AQUI?
	// 1. Evita cópias: mongo.Client é uma struct grande, passar por ponteiro
	//    evita copiar todos os dados a cada operação
	// 2. Compartilhamento: múltiplas partes do código podem usar o mesmo cliente
	//    sem criar cópias independentes
	// 3. Modificação: permite que métodos modifiquem o estado interno do cliente
	//
	// Exemplo prático:
	//   var x int = 10        // x é um valor
	//   var p *int = &x      // p é um ponteiro para x (armazena o endereço de x)
	//   *p = 20              // modifica x através do ponteiro (x agora é 20)
	client := mongo.NewClient(mongoURI)

	// defer garante que esta função seja executada quando main() terminar
	// Mesmo se houver um panic ou return antecipado, o defer sempre executa
	// Isso é essencial para limpar recursos (fechar conexões, arquivos, etc.)
	defer func() {
		if err := client.Disconnect(nil); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Database() retorna um ponteiro (*mongo.Database)
	// Todas as operações no banco usam este mesmo objeto compartilhado
	db := client.Database("userdb")

	// ============================================
	// INJEÇÃO DE DEPENDÊNCIAS
	// ============================================
	// Montamos a cadeia: Repository → UseCase → Handler
	// Cada camada recebe a anterior como dependência
	//
	// POR QUE FAZER ISSO?
	// 1. Testabilidade: podemos passar mocks (implementações falsas) para testar
	// 2. Flexibilidade: podemos trocar MongoDB por PostgreSQL sem mudar usecase/handler
	// 3. Desacoplamento: cada camada não conhece detalhes da implementação da outra
	//
	// O fluxo é: Handler usa UseCase, UseCase usa Repository, Repository usa MongoDB
	repo := repository.NewUserMongoRepository(db)
	uc := usecase.NewUserUseCase(repo)
	handler := httphandler.NewUserHandler(uc)

	// ============================================
	// CONFIGURAÇÃO DE ROTAS HTTP
	// ============================================
	// Chi é um router HTTP leve e rápido para Go
	// Router mapeia URLs para funções (handlers)
	r := chi.NewRouter()

	// Registra rota de healthcheck
	httphandler.RegisterHealth(r)

	// Registra rotas de usuários (CRUD)
	handler.RegisterRoutes(r)

	// Registra rotas do Swagger UI (documentação interativa)
	// Acesse: http://localhost:8080/swagger/index.html
	httphandler.RegisterSwagger(r)

	// ============================================
	// INICIALIZAÇÃO DO SERVIDOR
	// ============================================
	// ListenAndServe inicia um servidor HTTP que escuta na porta especificada
	// O segundo parâmetro é o handler (router) que processa as requisições
	//
	// IMPORTANTE: Esta função é BLOQUEANTE
	// Ela fica rodando indefinidamente até que o servidor seja encerrado
	// Por isso não há código depois dela - ela nunca retorna normalmente
	//
	// Em produção, considere:
	// - Adicionar timeouts (ReadTimeout, WriteTimeout)
	// - Configurar TLS/HTTPS
	// - Usar graceful shutdown (permitir requisições em andamento terminarem)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
