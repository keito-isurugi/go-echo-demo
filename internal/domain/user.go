package domain

type User struct {
	ID           int
	Name         string
	Email        string
	Password     string
	ProviderID   string
	ProviderName string
}

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}
