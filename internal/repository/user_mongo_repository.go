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

// userDoc represents the MongoDB document structure
type userDoc struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

// UserMongoRepository implements domain.UserRepository using MongoDB
type UserMongoRepository struct {
	collection *mongo.Collection
}

// NewUserMongoRepository creates a new MongoDB user repository
func NewUserMongoRepository(db *mongo.Database) domain.UserRepository {
	return &UserMongoRepository{
		collection: db.Collection("users"),
	}
}

// Create inserts a new user document
func (r *UserMongoRepository) Create(user *domain.User) error {
	// Cria um contexto com timeout para a operação de escrita.
	// Contexts evitam que operações pendentes fiquem indefinidamente travadas.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte a entidade de domínio para o formato de documento usado no MongoDB.
	doc := userDoc{
		Name:  user.Name,
		Email: user.Email,
	}

	// Insere o documento na collection. `InsertOne` usa o contexto com timeout.
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	// O driver retorna o _id gerado. Convertendo para ObjectID e depois para hex
	// para armazenar no campo ID da entidade de domínio.
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// GetByID retrieves a user by ID
func (r *UserMongoRepository) GetByID(id string) (*domain.User, error) {
	// Busca por ID exige converter o hex string para ObjectID do MongoDB.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Conversão do ID (linha importante): se a string não for um hex válido
	// de ObjectID, retornamos ErrNotFound para que a camada superior trate
	// como recurso inexistente.
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, usecase.ErrNotFound
	}

	var doc userDoc
	// Executa a consulta e decodifica o resultado no struct doc.
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, usecase.ErrNotFound
		}
		return nil, err
	}

	// Converte o documento do banco para a entidade de domínio.
	return &domain.User{
		ID:    doc.ID.Hex(),
		Name:  doc.Name,
		Email: doc.Email,
	}, nil
}

// List retrieves all users
func (r *UserMongoRepository) List() ([]*domain.User, error) {
	// Lista todos os documentos da collection `users`.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	// Garante fechamento do cursor ao final da função.
	defer cursor.Close(ctx)

	var users []*domain.User
	// Iteramos sobre o cursor e decodificamos cada documento.
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

	// Verifica se ocorreu erro durante a iteração do cursor.
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates a user document
func (r *UserMongoRepository) Update(user *domain.User) error {
	// Atualiza um documento pelo seu ObjectID.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Converte o ID recebido (hex) para ObjectID. Se inválido, retornamos
	// ErrNotFound para manter a API consistente (não vazamos detalhes do DB).
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return usecase.ErrNotFound
	}

	// Monta o update usando operador $set para alterar apenas campos informados.
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	// Executa a atualização por ID.
	result, err := r.collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}

	// Se nenhum documento foi correspondido, retorna ErrNotFound.
	if result.MatchedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

// Delete removes a user document
func (r *UserMongoRepository) Delete(id string) error {
	// Remove um documento pelo seu ID.
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

	if result.DeletedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}
