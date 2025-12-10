package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"user-api/internal/domain"
	"user-api/internal/usecase"
)

// ============================================
// ESTRUTURA PARA MONGODB
// ============================================
// userDoc é a estrutura usada para armazenar no MongoDB
// É diferente de domain.User porque o MongoDB usa ObjectID, não string
//
// SOBRE AS TAGS BSON:
// - `bson:"_id,omitempty"` significa:
//   * O campo ID no Go vira "_id" no MongoDB
//   * omitempty: se o campo estiver vazio, não inclui no documento
// - `bson:"name"` mapeia o campo Name para "name" no MongoDB
//
// POR QUE TER DUAS ESTRUTURAS (userDoc e domain.User)?
// - domain.User é a entidade do domínio (não conhece MongoDB)
// - userDoc é específica do MongoDB (usa ObjectID)
// - Fazemos conversão entre elas (isso é responsabilidade do repository)
// - Isso mantém o domínio independente do banco de dados
type userDoc struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`  // ObjectID é o tipo nativo do MongoDB
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

// ============================================
// REPOSITÓRIO MONGODB
// ============================================
// UserMongoRepository implementa domain.UserRepository usando MongoDB
//
// SOBRE collection *mongo.Collection:
// - collection é um ponteiro para a collection do MongoDB
// - Collection é como uma "tabela" no MongoDB
// - Todas as operações (insert, find, update, delete) usam esta collection
type UserMongoRepository struct {
	collection *mongo.Collection  // Ponteiro para a collection "users" do MongoDB
}

// NewUserMongoRepository cria um repositório MongoDB
//
// PARÂMETRO db *mongo.Database:
// - Recebe um ponteiro para o database
// - *mongo.Database significa "ponteiro para mongo.Database"
// - Usamos ponteiro para evitar copiar a struct (que pode ser grande)
//
// RETORNO &UserMongoRepository{...}:
// - O & cria um ponteiro para a struct UserMongoRepository
// - Retornamos ponteiro porque:
//   1. Evita cópia da struct (mais eficiente)
//   2. Permite que métodos modifiquem o estado interno (se necessário)
//   3. É padrão em Go retornar ponteiros de structs
//
// POR QUE RETORNAR domain.UserRepository (interface)?
// - Retornamos a interface, não o tipo concreto
// - Isso permite que o código que usa não dependa de MongoDB
// - Se mudarmos para PostgreSQL, só mudamos esta implementação
func NewUserMongoRepository(db *mongo.Database) domain.UserRepository {
	return &UserMongoRepository{
		collection: db.Collection("users"),  // Obtém a collection "users"
	}
}

// ============================================
// CREATE
// ============================================
// Create insere um novo usuário no MongoDB
// O ID é gerado automaticamente pelo MongoDB e depois convertido para string hex
//
// PARÂMETRO user *domain.User:
// - Recebe um ponteiro para poder MODIFICAR o campo ID
// - Quando o MongoDB gera o ID, precisamos colocá-lo de volta no user
// - Se recebêssemos domain.User (valor), modificaríamos apenas uma cópia
func (r *UserMongoRepository) Create(user *domain.User) error {
	// Context com timeout evita que a operação trave indefinidamente
	// Se o MongoDB estiver lento ou travado, após 5 segundos a operação cancela
	//
	// SOBRE CONTEXT:
	// - context.Background() cria um contexto vazio (raiz)
	// - WithTimeout adiciona um timeout de 5 segundos
	// - cancel() é uma função para cancelar manualmente (se necessário)
	// - defer cancel() garante que o contexto seja cancelado ao final
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte a entidade do domínio (domain.User) para o formato do MongoDB (userDoc)
	// Note: não incluímos o ID porque o MongoDB vai gerar automaticamente
	// O campo ID em userDoc tem tag `omitempty`, então será ignorado se vazio
	doc := userDoc{
		Name:  user.Name,
		Email: user.Email,
		// ID não é definido - MongoDB vai gerar automaticamente
	}

	// Insere o documento no MongoDB
	// InsertOne retorna um resultado com o ID gerado
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err  // Propaga o erro (ex: banco indisponível, conexão perdida)
	}

	// Pega o ID gerado pelo MongoDB e converte para string hexadecimal
	// 
	// SOBRE A CONVERSÃO:
	// - result.InsertedID é do tipo interface{} (tipo genérico)
	// - Fazemos type assertion: .(primitive.ObjectID) para converter
	// - .Hex() converte ObjectID para string hexadecimal
	//
	// POR QUE MODIFICAR user.ID AQUI?
	// - user é um ponteiro (*domain.User)
	// - Quando fazemos user.ID = ..., estamos modificando a struct apontada
	// - Essa modificação é visível para quem chamou Create()
	// - O usecase que chamou pode retornar o user com ID já populado
	//
	// Exemplo do que acontece:
	//   user := &domain.User{Name: "João"}  // user.ID = ""
	//   repo.Create(user)                    // Dentro: user.ID = "507f1f77..."
	//   // Agora user.ID tem valor mesmo fora do Create!
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// ============================================
// GET BY ID
// ============================================
// GetByID busca um usuário pelo ID
// Retorna um ponteiro (*domain.User) para evitar copiar a struct
func (r *UserMongoRepository) GetByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte a string hexadecimal para ObjectID do MongoDB
	// ObjectIDFromHex valida se a string é um hex válido de 24 caracteres
	//
	// Se o formato estiver inválido (ex: "abc", "123"), retorna erro
	// Nesse caso, retornamos ErrNotFound para manter a API consistente
	// (não vazamos detalhes técnicos do MongoDB para o usecase)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, usecase.ErrNotFound
	}

	// Declara uma variável do tipo userDoc (vazia)
	// O Decode vai preencher esta struct com os dados do MongoDB
	var doc userDoc
	
	// Busca o documento no MongoDB e decodifica no struct doc
	//
	// SOBRE bson.M{"_id": oid}:
	// - bson.M é um map[string]interface{} usado para queries MongoDB
	// - {"_id": oid} significa "buscar onde _id seja igual a oid"
	// - É equivalente a SQL: SELECT * FROM users WHERE _id = oid
	//
	// SOBRE .Decode(&doc):
	// - Decode converte o documento BSON do MongoDB para a struct Go
	// - O & passa um ponteiro para doc, permitindo que Decode preencha os campos
	// - Se não passar ponteiro, Decode não conseguiria modificar doc
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		// Se não encontrar documento, retorna erro específico
		if err == mongo.ErrNoDocuments {
			return nil, usecase.ErrNotFound
		}
		// Outros erros (ex: conexão perdida) são propagados
		return nil, err
	}

	// Converte de volta para a entidade do domínio
	// Retornamos um ponteiro usando & para evitar cópia
	//
	// POR QUE RETORNAR PONTEIRO?
	// - domain.User pode crescer (adicionar mais campos)
	// - Retornar ponteiro é mais eficiente (não copia a struct)
	// - Permite que o chamador modifique se necessário (embora não façamos isso)
	return &domain.User{
		ID:    doc.ID.Hex(),      // Converte ObjectID para string hex
		Name:  doc.Name,
		Email: doc.Email,
	}, nil
}

// ============================================
// LIST
// ============================================
// List retorna todos os usuários
// Retorna []*domain.User (slice de ponteiros) - mais eficiente que []domain.User
func (r *UserMongoRepository) List() ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Busca todos os documentos
	// bson.M{} significa "sem filtro" (equivalente a SELECT * FROM users)
	// Find retorna um Cursor, que é um iterador sobre os resultados
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	// Garante que o cursor seja fechado ao final (libera recursos)
	defer cursor.Close(ctx)

	// Cria um slice vazio de ponteiros para domain.User
	// []*domain.User significa "slice de ponteiros para domain.User"
	//
	// POR QUE SLICE DE PONTEIROS?
	// - Se fosse []domain.User, cada append copiaria a struct inteira
	// - Com []*domain.User, apenas copiamos o ponteiro (8 bytes) em vez da struct
	// - Mais eficiente, especialmente com muitos usuários
	var users []*domain.User
	
	// Itera sobre o cursor convertendo cada documento
	// cursor.Next() retorna true enquanto houver mais documentos
	for cursor.Next(ctx) {
		var doc userDoc
		
		// Decode converte o documento atual do cursor para a struct doc
		// O & passa ponteiro para doc, permitindo que Decode preencha os campos
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		
		// Cria um novo domain.User e adiciona ao slice
		// O & cria um ponteiro para a struct criada
		// append adiciona o ponteiro ao slice (não copia a struct)
		users = append(users, &domain.User{
			ID:    doc.ID.Hex(),
			Name:  doc.Name,
			Email: doc.Email,
		})
	}

	// Verifica se houve erro durante a iteração do cursor
	// Pode acontecer se a conexão cair no meio da leitura
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ============================================
// UPDATE
// ============================================
// Update atualiza um usuário existente
// Recebe *domain.User (ponteiro) com os campos já modificados pelo usecase
func (r *UserMongoRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte o ID (string hex) para ObjectID do MongoDB
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return usecase.ErrNotFound
	}

	// Monta a operação de update usando o operador $set
	// $set atualiza apenas os campos especificados, mantendo os outros intactos
	//
	// SOBRE $set:
	// - É um operador de update do MongoDB
	// - Atualiza apenas os campos listados
	// - Se o campo não existir, cria; se existir, atualiza
	//
	// Exemplo: se o documento tiver {_id: ..., name: "João", email: "joao@email.com", age: 30}
	// e fizermos $set: {name: "Maria"}, o resultado será:
	// {_id: ..., name: "Maria", email: "joao@email.com", age: 30}
	// (email e age permanecem inalterados)
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	// Executa a atualização no MongoDB
	// UpdateByID atualiza o documento com o _id especificado
	result, err := r.collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}

	// Verifica se algum documento foi encontrado e atualizado
	// MatchedCount = 0 significa que o ID não existe no banco
	if result.MatchedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

// ============================================
// DELETE
// ============================================
// Delete remove um usuário
func (r *UserMongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte o ID para ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return usecase.ErrNotFound
	}

	// Remove o documento do MongoDB
	// DeleteOne remove apenas um documento (o primeiro que encontrar)
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	// Verifica se algum documento foi deletado
	// DeletedCount = 0 significa que o ID não existe no banco
	if result.DeletedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}
