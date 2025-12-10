package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient cria e conecta um cliente MongoDB
// Retorna um ponteiro (*mongo.Client) para que todas as operações usem
// o mesmo cliente compartilhado, evitando cópias desnecessárias
func NewClient(uri string) *mongo.Client {
	// Context com timeout evita que a conexão trave indefinidamente
	// se o MongoDB não estiver disponível
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configura as opções de conexão usando a URI
	clientOptions := options.Client().ApplyURI(uri)

	// Tenta conectar ao MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Faz um ping para verificar se a conexão está realmente funcionando
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Retorna o cliente pronto para uso
	// Importante: quem chamar esta função deve fazer client.Disconnect() ao final
	return client
}
