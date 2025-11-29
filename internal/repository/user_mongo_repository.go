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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := userDoc{
		Name:  user.Name,
		Email: user.Email,
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// GetByID retrieves a user by ID
func (r *UserMongoRepository) GetByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, usecase.ErrNotFound
	}

	var doc userDoc
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, usecase.ErrNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:    doc.ID.Hex(),
		Name:  doc.Name,
		Email: doc.Email,
	}, nil
}

// List retrieves all users
func (r *UserMongoRepository) List() ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
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

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates a user document
func (r *UserMongoRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return usecase.ErrNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	result, err := r.collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

// Delete removes a user document
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

	if result.DeletedCount == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

