package usecase

import (
	"errors"
	"strings"

	"user-api/internal/domain"
)

// Erros customizados usados pela camada de usecase
var (
	ErrInvalidEmail = errors.New("invalid email")  // Email sem '@'
	ErrNotFound     = errors.New("user not found")  // Usuário não encontrado
)

// userUseCase implementa a lógica de negócio
// Contém validações e orquestra as chamadas ao repositório
type userUseCase struct {
	repo domain.UserRepository
}

// NewUserUseCase cria um novo usecase recebendo o repositório como dependência
// Isso permite trocar a implementação (MongoDB, memória para testes, etc.)
func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{repo: repo}
}

// CreateUser valida o email e cria um novo usuário
// O repositório vai popular o campo ID quando persistir no banco
func (uc *userUseCase) CreateUser(name, email string) (*domain.User, error) {
	// Validação básica: email deve conter '@'
	// Em produção, use uma biblioteca de validação mais robusta
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	// Cria a entidade. O & cria um ponteiro, permitindo que o repositório
	// modifique o ID diretamente na mesma instância
	user := &domain.User{
		Name:  name,
		Email: email,
	}

	// Persiste no banco. Se der erro, propaga para o handler
	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser busca um usuário por ID
func (uc *userUseCase) GetUser(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

// ListUsers retorna todos os usuários
func (uc *userUseCase) ListUsers() ([]*domain.User, error) {
	return uc.repo.List()
}

// UpdateUser atualiza os campos de um usuário existente
// Só atualiza campos que foram informados (não vazios)
func (uc *userUseCase) UpdateUser(id, name, email string) (*domain.User, error) {
	// Primeiro busca o usuário atual
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	// Atualiza apenas os campos informados
	if name != "" {
		user.Name = name
	}

	if email != "" {
		// Valida o novo email se foi informado
		if !strings.Contains(email, "@") {
			return nil, ErrInvalidEmail
		}
		user.Email = email
	}

	// Salva as alterações
	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser remove um usuário
func (uc *userUseCase) DeleteUser(id string) error {
	return uc.repo.Delete(id)
}
