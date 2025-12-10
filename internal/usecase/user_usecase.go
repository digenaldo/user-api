package usecase

import (
	"errors"
	"strings"

	"user-api/internal/domain"
)

// Variáveis de erro exportadas
var (
	// ErrInvalidEmail é retornado quando o e-mail não contém o caractere '@'
	// Utilizamos um erro simples para exemplificar validação — em produção
	// você pode aplicar validações mais robustas (regex, verificação MX, etc.).
	ErrInvalidEmail = errors.New("invalid email")
	// ErrNotFound é retornado quando o recurso não é encontrado no repositório.
	ErrNotFound = errors.New("user not found")
)

// userUseCase implementa a interface domain.UserUseCase e contém a lógica
// de negócio (validações, regras e orquestração entre repositório e handlers).
type userUseCase struct {
	repo domain.UserRepository
}

// NewUserUseCase cria uma instância do usecase a partir de um repositório.
// Recebemos a abstração `domain.UserRepository` para permitir testes e
// troca de implementação (por exemplo, memória, MongoDB, etc.).
func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{repo: repo}
}

// CreateUser cria um novo usuário após validações básicas.
// Explicação passo-a-passo (trecho importante):
//   - Verificamos se o e-mail contém '@' (validação mínima).
//   - Construímos a entidade `domain.User` e delegamos para o repositório.
//   - O repositório popula o ID (quando aplicável) e pode retornar erros
//     relacionados à persistência.
func (uc *userUseCase) CreateUser(name, email string) (*domain.User, error) {
	// Validação simples de e-mail (linha importante)
	if !strings.Contains(email, "@") {
		// Retornamos um erro de negócio específico para o caso de uso
		return nil, ErrInvalidEmail
	}

	// Construímos a entidade de domínio. Note que aqui ainda não há ID,
	// pois o repositório (MongoDB) deverá atribuí-lo ao persistir.
	//
	// Observação sobre `&` e `*`:
	// - O operador `&` antes de `domain.User{...}` cria uma struct e
	//   retorna seu endereço (um *domain.User). Ou seja `user` é um
	//   ponteiro para a entidade.
	// - Usamos ponteiro para que, ao passá-lo para o repositório
	//   (`repo.Create(user)`), o repositório possa preencher/alterar
	//   campos (por exemplo `ID`) diretamente na mesma instância —
	//   evitando cópias desnecessárias e permitindo que o chamador
	//   observe as modificações.
	user := &domain.User{
		Name:  name,
		Email: email,
	}

	// Persiste o usuário via repositório. Erros de banco são propagados
	// para o chamador — o handler traduz em resposta HTTP adequada.
	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser busca um usuário por ID delegando diretamente ao repositório.
// A validação do formato do ID é tratada pelo repositório (ex: conversão
// para ObjectID do MongoDB).
func (uc *userUseCase) GetUser(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

// ListUsers retorna todos os usuários; apenas repassa a chamada ao repo.
func (uc *userUseCase) ListUsers() ([]*domain.User, error) {
	return uc.repo.List()
}

// UpdateUser atualiza os campos do usuário.
// Fluxo principal:
// - Busca o usuário atual (GetByID)
// - Se não existir, retorna ErrNotFound
// - Atualiza apenas campos não vazios
// - Valida o e-mail quando informado
// - Persiste a atualização via repositório
func (uc *userUseCase) UpdateUser(id, name, email string) (*domain.User, error) {
	// Recupera o usuário atual
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Segurança: se o repositório puder retornar nil (sem erro), tratamos
	// como recurso não encontrado.
	if user == nil {
		return nil, ErrNotFound
	}

	// Atualiza campos se fornecidos
	if name != "" {
		user.Name = name
	}

	if email != "" {
		// Validação mínima do e-mail (trecho importante)
		if !strings.Contains(email, "@") {
			return nil, ErrInvalidEmail
		}
		user.Email = email
	}

	// Persiste a atualização
	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser delega a remoção do usuário ao repositório.
func (uc *userUseCase) DeleteUser(id string) error {
	return uc.repo.Delete(id)
}
