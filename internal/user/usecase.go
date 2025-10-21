package user

import (
	"go-clean-api/internal/user/domain"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
}

type userUsecase struct {
	userRepo UserRepository
}

func NewUserUsecase(repo UserRepository) UserUsecase {
	return &userUsecase{userRepo: repo}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	return u.userRepo.Create(user)
}

func (u *userUsecase) GetAllUsers() ([]domain.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) GetUserByID(id uint) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) UpdateUser(user *domain.User) error {
	return u.userRepo.Update(user)
}

func (u *userUsecase) DeleteUser(id uint) error {
	return u.userRepo.Delete(id)
}
