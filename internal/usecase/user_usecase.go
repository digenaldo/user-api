package usecase

import (
	"errors"
	"strings"

	"user-api/internal/domain"
)

var (
	// ErrInvalidEmail is returned when email format is invalid
	ErrInvalidEmail = errors.New("invalid email")
	// ErrNotFound is returned when user is not found
	ErrNotFound = errors.New("user not found")
)

// userUseCase implements domain.UserUseCase
type userUseCase struct {
	repo domain.UserRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{repo: repo}
}

// CreateUser creates a new user
func (uc *userUseCase) CreateUser(name, email string) (*domain.User, error) {
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	user := &domain.User{
		Name:  name,
		Email: email,
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (uc *userUseCase) GetUser(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

// ListUsers retrieves all users
func (uc *userUseCase) ListUsers() ([]*domain.User, error) {
	return uc.repo.List()
}

// UpdateUser updates a user
func (uc *userUseCase) UpdateUser(id, name, email string) (*domain.User, error) {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	if name != "" {
		user.Name = name
	}

	if email != "" {
		if !strings.Contains(email, "@") {
			return nil, ErrInvalidEmail
		}
		user.Email = email
	}

	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (uc *userUseCase) DeleteUser(id string) error {
	return uc.repo.Delete(id)
}
