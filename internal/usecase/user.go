package usecase

import (
	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/repository"
)

type UserUsecase interface {
	GetUsers() ([]domain.User, error)
	GetUser(id int) (*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(id int) error
}

type userUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{repo: repo}
}

func (u *userUsecaseImpl) GetUsers() ([]domain.User, error) {
	return u.repo.FindAll()
}

func (u *userUsecaseImpl) GetUser(id int) (*domain.User, error) {
	return u.repo.FindByID(id)
}

func (u *userUsecaseImpl) CreateUser(user *domain.User) error {
	return u.repo.Create(user)
}

func (u *userUsecaseImpl) UpdateUser(user *domain.User) error {
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) DeleteUser(id int) error {
	return u.repo.Delete(id)
}
