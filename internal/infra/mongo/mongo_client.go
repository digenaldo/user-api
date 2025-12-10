package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ============================================
// CLIENTE MONGODB
// ============================================
// NewClient cria e conecta um cliente MongoDB
//
// RETORNA *mongo.Client (ponteiro):
// - Retorna um ponteiro para que todas as operações usem o mesmo cliente compartilhado
// - Evita cópias desnecessárias da struct mongo.Client (que pode ser grande)
// - Permite que múltiplas partes do código usem o mesmo cliente
//
// SOBRE PONTEIROS EM GO:
// - *T significa "ponteiro para tipo T"
// - &x cria um ponteiro para x (pega o endereço de x)
// - *p acessa o valor apontado por p (desreferencia o ponteiro)
//
// Exemplo prático:
//   var x int = 10
//   var p *int = &x    // p aponta para x
//   *p = 20            // modifica x através do ponteiro
//   // x agora é 20
func NewClient(uri string) *mongo.Client {
	// Context com timeout evita que a conexão trave indefinidamente
	// Se o MongoDB não estiver disponível, após 10 segundos a operação cancela
	//
	// SOBRE CONTEXT:
	// - context.Background() cria um contexto raiz (vazio)
	// - WithTimeout adiciona um timeout de 10 segundos
	// - cancel() é uma função para cancelar manualmente (se necessário)
	// - defer cancel() garante que o contexto seja cancelado ao final da função
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configura as opções de conexão usando a URI fornecida
	// A URI tem formato: mongodb://usuario:senha@host:porta/database?opcoes
	clientOptions := options.Client().ApplyURI(uri)

	// Tenta conectar ao MongoDB
	// mongo.Connect retorna (*mongo.Client, error)
	// Se falhar (ex: URI inválida, servidor inacessível), retorna erro
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		// log.Fatalf encerra a aplicação imediatamente
		// Usamos Fatal porque sem MongoDB a aplicação não funciona
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Faz um ping para verificar se a conexão está realmente funcionando
	// Só conectar não garante que o servidor está respondendo
	// O ping confirma que conseguimos se comunicar com o MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Retorna o cliente pronto para uso
	// IMPORTANTE: quem chamar esta função deve fazer client.Disconnect() ao final
	// Isso libera os recursos de conexão (sockets, goroutines, etc.)
	return client
}
