package domain

// User representa um usuário na aplicação
// O ID é uma string hexadecimal do ObjectID do MongoDB (ex: "507f1f77bcf86cd799439011")
// A validação do email (deve conter '@') é feita na camada de usecase
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRepository define o contrato para persistência de usuários
// Usamos interface para poder trocar a implementação (MongoDB, PostgreSQL, memória, etc.)
// Isso facilita testes e permite mudar de banco sem alterar o resto do código
type UserRepository interface {
	Create(user *User) error      // Cria um usuário e popula o ID
	GetByID(id string) (*User, error) // Busca por ID, retorna erro se não encontrar
	List() ([]*User, error)       // Lista todos os usuários
	Update(user *User) error      // Atualiza um usuário existente
	Delete(id string) error       // Remove um usuário pelo ID
}

// UserUseCase define o contrato da camada de lógica de negócio
// Os handlers chamam esses métodos para executar as regras do domínio
type UserUseCase interface {
	CreateUser(name, email string) (*User, error)  // Valida e cria usuário
	GetUser(id string) (*User, error)              // Busca por ID
	ListUsers() ([]*User, error)                   // Lista todos
	UpdateUser(id, name, email string) (*User, error) // Atualiza campos
	DeleteUser(id string) error                    // Remove usuário
}
