package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient creates a new MongoDB client
func NewClient(uri string) *mongo.Client {
	// Cria um contexto com timeout para as operações iniciais de conexão.
	// Usamos um timeout para evitar que a inicialização trave indefinidamente
	// caso o banco não esteja alcançável.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configura opções do cliente usando a URI fornecida.
	clientOptions := options.Client().ApplyURI(uri)

	// Conecta ao MongoDB. A função pode falhar se a URI estiver incorreta
	// ou se o serviço não estiver disponível.
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		// Falha na conexão é considerada erro crítico na inicialização.
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verifica a conectividade com um `Ping`. Passamos o mesmo contexto
	// com timeout para garantir que a verificação não demore demais.
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Retornamos o cliente conectado. Quem chamar esta função deve garantir
	// que `client.Disconnect(ctx)` seja chamado ao finalizar (ver cmd/api/main.go).
	return client
}
