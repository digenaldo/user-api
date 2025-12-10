package usecase

import (
	"errors"
	"strings"

	"user-api/internal/domain"
)

// ============================================
// ERROS CUSTOMIZADOS
// ============================================
// Erros customizados permitem identificar tipos específicos de erro
// Isso é útil para o handler decidir qual status HTTP retornar
//
// POR QUE ERRORS.NEW()?
// - Cria um erro simples com uma mensagem
// - Podemos comparar erros usando == (err == ErrInvalidEmail)
// - Mais simples que criar structs complexas para erros
var (
	ErrInvalidEmail = errors.New("invalid email")  // Email sem '@'
	ErrNotFound     = errors.New("user not found")  // Usuário não encontrado
)

// ============================================
// IMPLEMENTAÇÃO DO USECASE
// ============================================
// userUseCase implementa a lógica de negócio
// Contém validações e orquestra as chamadas ao repositório
//
// SOBRE O RECEPTOR (uc *userUseCase):
// - O receptor é o "this" ou "self" de outras linguagens
// - *userUseCase significa que o receptor é um PONTEIRO
// - Isso permite que métodos modifiquem o estado interno (se houver)
// - É uma prática comum em Go usar ponteiros como receptores
type userUseCase struct {
	repo domain.UserRepository  // Dependência: o repositório que vamos usar
}

// NewUserUseCase cria um novo usecase recebendo o repositório como dependência
// Isso permite trocar a implementação (MongoDB, memória para testes, etc.)
//
// POR QUE RETORNAR &userUseCase{...} (ponteiro)?
// - Retornamos um ponteiro para que todas as chamadas usem a mesma instância
// - Se retornássemos userUseCase (valor), cada chamada criaria uma cópia
// - O & cria um ponteiro para a struct criada
func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{repo: repo}
}

// ============================================
// CREATE USER
// ============================================
// CreateUser valida o email e cria um novo usuário
// O repositório vai popular o campo ID quando persistir no banco
func (uc *userUseCase) CreateUser(name, email string) (*domain.User, error) {
	// Validação básica: email deve conter '@'
	// Em produção, use uma biblioteca de validação mais robusta (ex: validator)
	// Poderia validar: formato correto, domínio válido, não estar em blacklist, etc.
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	// Cria a entidade usando o operador & (address-of)
	// &domain.User{...} cria uma struct e retorna um PONTEIRO para ela
	//
	// POR QUE PONTEIRO AQUI?
	// 1. O repositório precisa MODIFICAR o campo ID após salvar
	// 2. Se passássemos domain.User (valor), o repositório modificaria uma CÓPIA
	// 3. Com ponteiro, o repositório modifica a MESMA instância
	// 4. Quando retornamos user, ele já terá o ID populado
	//
	// Exemplo do que acontece:
	//   user := &domain.User{Name: "João", Email: "joao@email.com"}
	//   // user.ID ainda está vazio ""
	//   repo.Create(user)  // Passa o ponteiro
	//   // Dentro do Create, fazemos: user.ID = "507f1f77..."
	//   // Como user é ponteiro, essa mudança é visível aqui também!
	//   return user  // user.ID agora tem valor
	user := &domain.User{
		Name:  name,
		Email: email,
		// ID ainda está vazio - será populado pelo repositório
	}

	// Persiste no banco através do repositório
	// Se der erro (ex: banco indisponível), propaga para o handler
	// O handler decide como tratar (retornar 500, 503, etc.)
	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	// Retorna o usuário criado (agora com ID populado)
	// Como user é um ponteiro, retornamos o mesmo ponteiro
	return user, nil
}

// ============================================
// GET USER
// ============================================
// GetUser busca um usuário por ID
// Apenas repassa a chamada para o repositório
// A lógica de negócio aqui é mínima - poderia adicionar cache, logging, etc.
func (uc *userUseCase) GetUser(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

// ============================================
// LIST USERS
// ============================================
// ListUsers retorna todos os usuários
// Em uma aplicação real, poderia adicionar:
// - Paginação (limite, offset)
// - Filtros (buscar por nome, email)
// - Ordenação (por nome, data de criação)
func (uc *userUseCase) ListUsers() ([]*domain.User, error) {
	return uc.repo.List()
}

// ============================================
// UPDATE USER
// ============================================
// UpdateUser atualiza os campos de um usuário existente
// Só atualiza campos que foram informados (não vazios)
//
// FLUXO:
// 1. Busca o usuário atual no banco
// 2. Verifica se existe
// 3. Atualiza apenas campos não vazios
// 4. Valida email se foi informado
// 5. Salva as alterações
func (uc *userUseCase) UpdateUser(id, name, email string) (*domain.User, error) {
	// Primeiro busca o usuário atual
	// GetByID retorna (*User, error)
	// Se não encontrar, retorna (nil, ErrNotFound)
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verificação de segurança: se o repositório retornar nil sem erro
	// (não deveria acontecer, mas é bom prevenir)
	if user == nil {
		return nil, ErrNotFound
	}

	// Atualiza apenas os campos informados (não vazios)
	// Isso permite atualizar apenas name OU apenas email
	//
	// SOBRE MODIFICAR user DIRETAMENTE:
	// - user é um ponteiro (*domain.User) retornado pelo repositório
	// - Quando modificamos user.Name, estamos modificando a struct apontada
	// - Essa modificação será persistida quando chamarmos repo.Update(user)
	// - Não precisamos criar uma nova struct - modificamos a existente
	if name != "" {
		user.Name = name
	}

	if email != "" {
		// Valida o novo email se foi informado
		// Mesma validação do CreateUser
		if !strings.Contains(email, "@") {
			return nil, ErrInvalidEmail
		}
		user.Email = email
	}

	// Salva as alterações no banco
	// O repositório recebe o ponteiro user com os campos já modificados
	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	// Retorna o usuário atualizado
	// Como user é um ponteiro, retornamos o mesmo ponteiro (mesma instância)
	return user, nil
}

// ============================================
// DELETE USER
// ============================================
// DeleteUser remove um usuário
// Apenas repassa para o repositório
// Poderia adicionar: soft delete, verificar dependências, etc.
func (uc *userUseCase) DeleteUser(id string) error {
	return uc.repo.Delete(id)
}
