package domain

// User represents a user entity
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	List() ([]*User, error)
	Update(user *User) error
	Delete(id string) error
}

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	CreateUser(name, email string) (*User, error)
	GetUser(id string) (*User, error)
	ListUsers() ([]*User, error)
	UpdateUser(id, name, email string) (*User, error)
	DeleteUser(id string) error
}

