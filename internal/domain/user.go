package domain

// ============================================
// ENTIDADE DE DOMÍNIO
// ============================================
// User representa um usuário na aplicação
// Esta é a entidade central do nosso domínio - ela não depende de nada externo
//
// SOBRE AS TAGS JSON:
// - `json:"id"` significa que ao serializar para JSON, o campo ID vira "id"
// - Isso permite ter nomes diferentes em Go (maiúsculo) e JSON (minúsculo)
// - Quando recebemos JSON, o Go automaticamente mapeia usando essas tags
//
// O ID é uma string hexadecimal do ObjectID do MongoDB
// Exemplo: "507f1f77bcf86cd799439011"
// A validação do email (deve conter '@') é feita na camada de usecase
type User struct {
	ID    string `json:"id"`    // Identificador único (hex do ObjectID do MongoDB)
	Name  string `json:"name"`  // Nome completo do usuário
	Email string `json:"email"`  // Email (deve conter '@')
}

// ============================================
// INTERFACE DO REPOSITORY
// ============================================
// UserRepository define o contrato (interface) para persistência de usuários
//
// POR QUE USAR INTERFACE?
// 1. Abstração: define O QUE fazer, não COMO fazer
// 2. Flexibilidade: podemos ter implementações diferentes (MongoDB, PostgreSQL, memória)
// 3. Testabilidade: podemos criar um "mock" (implementação falsa) para testes
// 4. Desacoplamento: o usecase não precisa saber que usamos MongoDB
//
// SOBRE PONTEIROS NOS PARÂMETROS:
// - Create(user *User): recebe ponteiro para poder MODIFICAR o user (popular o ID)
// - GetByID retorna (*User, error): retorna ponteiro para evitar cópia da struct
// - List retorna ([]*User, error): slice de ponteiros - mais eficiente para muitos itens
// - Update(user *User): recebe ponteiro para modificar os campos
//
// Quando usar ponteiro vs valor?
// - Use ponteiro (*T) quando: struct é grande, precisa modificar, quer compartilhar
// - Use valor (T) quando: struct é pequena, não precisa modificar, quer cópia independente
type UserRepository interface {
	// Create persiste um novo usuário
	// Recebe *User (ponteiro) para poder popular o campo ID após salvar
	// O repositório modifica o user.ID diretamente na mesma instância
	Create(user *User) error
	
	// GetByID busca um usuário pelo ID
	// Retorna *User (ponteiro) para evitar copiar a struct
	// Se não encontrar, retorna erro (não retorna nil sem erro)
	GetByID(id string) (*User, error)
	
	// List retorna todos os usuários
	// Retorna []*User (slice de ponteiros) - mais eficiente que []User
	// Cada elemento do slice é um ponteiro para uma struct User
	List() ([]*User, error)
	
	// Update atualiza um usuário existente
	// Recebe *User (ponteiro) com os campos já modificados
	// O repositório apenas persiste as alterações
	Update(user *User) error
	
	// Delete remove um usuário pelo ID
	// Retorna apenas error (não precisa retornar o usuário deletado)
	Delete(id string) error
}

// ============================================
// INTERFACE DO USECASE
// ============================================
// UserUseCase define o contrato da camada de lógica de negócio
// Os handlers HTTP chamam esses métodos para executar as regras do domínio
//
// DIFERENÇA ENTRE REPOSITORY E USECASE:
// - Repository: cuida de COMO salvar/buscar dados (detalhes técnicos)
// - UseCase: cuida de O QUE fazer com os dados (regras de negócio, validações)
//
// Exemplo: Repository sabe converter ObjectID, UseCase sabe validar email
type UserUseCase interface {
	// CreateUser valida os dados e cria um novo usuário
	// Retorna *User (ponteiro) com o usuário criado (incluindo o ID gerado)
	CreateUser(name, email string) (*User, error)
	
	// GetUser busca um usuário pelo ID
	// Retorna *User (ponteiro) ou erro se não encontrar
	GetUser(id string) (*User, error)
	
	// ListUsers retorna todos os usuários cadastrados
	// Retorna []*User (slice de ponteiros)
	ListUsers() ([]*User, error)
	
	// UpdateUser atualiza os campos de um usuário existente
	// Recebe id e os novos valores (name e email podem ser vazios)
	// Retorna *User (ponteiro) com os dados atualizados
	UpdateUser(id, name, email string) (*User, error)
	
	// DeleteUser remove um usuário pelo ID
	// Retorna apenas error (não precisa retornar o usuário deletado)
	DeleteUser(id string) error
}
