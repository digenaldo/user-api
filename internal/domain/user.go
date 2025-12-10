package domain

// User é a entidade de domínio que representa um usuário na aplicação.
// Comentários didáticos por campo:
//   - ID: string que representa o identificador único do usuário. No caso
//     desta aplicação usamos o hex do ObjectID do MongoDB (ex: "507f1f77...").
//   - Name: nome completo do usuário.
//   - Email: email do usuário; espera-se que contenha '@' (validação feita
//     na camada de usecase).
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRepository define a interface que qualquer repositório de usuários
// deve implementar. A ideia é depender de abstrações (interfaces) e não
// de implementações concretas — isso facilita testes e troca de banco.
// Métodos explicados:
//   - Create: persiste um novo usuário; pode popular o campo ID.
//   - GetByID: retorna um usuário pelo seu ID (string). Deve devolver um
//     erro ou nil quando não encontrado (a camada acima traduz para 404).
//   - List: retorna todos os usuários.
//   - Update: atualiza um usuário existente (recebe a entidade completa).
//   - Delete: remove um usuário pelo ID.
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	List() ([]*User, error)
	Update(user *User) error
	Delete(id string) error
}

// UserUseCase define a interface da camada de lógica de negócio (usecases).
// Os handlers HTTP chamam esses métodos para executar regras do domínio.
// Explicação dos métodos:
// - CreateUser: valida dados e cria um usuário retornando a entidade criada.
// - GetUser: recupera por ID.
// - ListUsers: lista todos.
// - UpdateUser: atualiza campos do usuário (recebe id e valores a alterar).
// - DeleteUser: remove usuário por id.
type UserUseCase interface {
	CreateUser(name, email string) (*User, error)
	GetUser(id string) (*User, error)
	ListUsers() ([]*User, error)
	UpdateUser(id, name, email string) (*User, error)
	DeleteUser(id string) error
}
