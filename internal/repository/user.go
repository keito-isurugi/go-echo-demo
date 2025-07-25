package repository

import "go-echo-demo/internal/domain"

type UserRepository interface {
	FindAll() ([]domain.User, error)
	FindByID(id int) (*domain.User, error)
	GetByID(id int) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id int) error
}
