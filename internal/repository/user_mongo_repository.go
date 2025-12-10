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

// userDoc é a estrutura usada para armazenar no MongoDB
// Usa primitive.ObjectID que é o tipo nativo do MongoDB para IDs
type userDoc struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

// UserMongoRepository implementa domain.UserRepository usando MongoDB
type UserMongoRepository struct {
	collection *mongo.Collection
}

// NewUserMongoRepository cria um repositório MongoDB
// Recebe o database como ponteiro para evitar cópias e garantir que todas
// as operações usem o mesmo cliente compartilhado
func NewUserMongoRepository(db *mongo.Database) domain.UserRepository {
	return &UserMongoRepository{
		collection: db.Collection("users"),
	}
}

// Create insere um novo usuário no MongoDB
// O ID é gerado automaticamente pelo MongoDB e depois convertido para string hex
func (r *UserMongoRepository) Create(user *domain.User) error {
	// Context com timeout evita que a operação trave indefinidamente
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte a entidade do domínio para o formato do MongoDB
	doc := userDoc{
		Name:  user.Name,
		Email: user.Email,
	}

	// Insere no banco
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	// Pega o ID gerado, converte para ObjectID e depois para string hex
	// Como user é um ponteiro, essa alteração é visível para quem chamou
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// GetByID busca um usuário pelo ID
func (r *UserMongoRepository) GetByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte a string hex para ObjectID do MongoDB
	// Se o formato estiver inválido, retorna erro de não encontrado
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, usecase.ErrNotFound
	}

	var doc userDoc
	// Busca o documento e decodifica no struct
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, usecase.ErrNotFound
		}
		return nil, err
	}

	// Converte de volta para a entidade do domínio
	return &domain.User{
		ID:    doc.ID.Hex(),
		Name:  doc.Name,
		Email: doc.Email,
	}, nil
}

// List retorna todos os usuários
func (r *UserMongoRepository) List() ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Busca todos os documentos (bson.M{} significa "sem filtro")
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	// Itera sobre o cursor convertendo cada documento
	for cursor.Next(ctx) {
		var doc userDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		users = append(users, &domain.User{
			ID:    doc.ID.Hex(),
			Name:  doc.Name,
			Email: doc.Email,
		})
	}

	// Verifica se houve erro durante a iteração
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update atualiza um usuário existente
func (r *UserMongoRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte o ID para ObjectID
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return usecase.ErrNotFound
	}

	// Usa $set para atualizar apenas os campos informados
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	// Executa a atualização
	result, err := r.collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}

	// Se nenhum documento foi encontrado, retorna erro
	if result.MatchedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

// Delete remove um usuário
func (r *UserMongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return usecase.ErrNotFound
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	// Se nenhum documento foi deletado, retorna erro
	if result.DeletedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}
