package domain

import (
	"go-clean-api/entity"

	"github.com/sirupsen/logrus"
)

type UserUsecase interface {
	CreateUser(log *logrus.Entry, user *entity.User) error
	GetAllUsers(log *logrus.Entry) ([]entity.User, error)
	GetUserByID(log *logrus.Entry, id uint) (*entity.User, error)
	UpdateUser(log *logrus.Entry, user *entity.User) error
	DeleteUser(log *logrus.Entry, id uint) error

	GetUserByEmail(log *logrus.Entry, email string) (*entity.User, error)
	ValidateUserCredentials(log *logrus.Entry, email, password string) (*entity.User, error)

	Login(log *logrus.Entry, email, password string) (string, error)
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
