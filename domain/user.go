package domain

import (
	"go-clean-api/entity"
)

type UserUsecase interface {
	CreateUser(user *entity.User) error
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uint) error

	GetUserByEmail(email string) (*entity.User, error)
	ValidateUserCredentials(email, password string) (*entity.User, error)

	Login(email, password string) (string, error)
}

type UserRepository interface {
	Create(user *entity.User) error
	FindAll() ([]entity.User, error)
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
}

// type UserCacheRepository interface {
// 	SetUserCache(userID uint, user *entity.User) error
// 	GetUserCache(userID uint) (*entity.User, error)
// 	DeleteUserCache(userID uint) error
// 	SetUserSession(sessionID string, userID uint) error
// 	GetUserSession(sessionID string) (uint, error)
// }

// type UserMessageRepository interface {
// 	PublishUserCreated(user *entity.User) error
// 	PublishUserUpdated(user *entity.User) error
// 	PublishUserDeleted(userID uint) error
// 	SubscribeUserEvents() error
// }
