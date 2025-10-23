// user/usecase/usecase.go
package usecase

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type userUsecase struct {
	userRepo    domain.UserRepository
	cacheRepo   domain.UserCacheRepository
	messageRepo domain.UserMessageRepository
}

func NewUserUsecase(
	userRepo domain.UserRepository,
	cacheRepo domain.UserCacheRepository,
	messageRepo domain.UserMessageRepository,
) domain.UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		cacheRepo:   cacheRepo,
		messageRepo: messageRepo,
	}
}

func (u *userUsecase) CreateUser(user *entity.User) error {
	// 1. Create in database
	if err := u.userRepo.Create(user); err != nil {
		return err
	}

	// 2. Cache the user
	go func() {
		u.cacheRepo.SetUserCache(uint(user.ID), user)
	}()

	// 3. Publish event
	go func() {
		u.messageRepo.PublishUserCreated(user)
	}()

	return nil
}

func (u *userUsecase) GetAllUsers() ([]entity.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) GetUserByID(id uint) (*entity.User, error) {
	// 1. Try cache first
	if user, err := u.cacheRepo.GetUserCache(id); err == nil && user != nil {
		return user, nil
	}

	// 2. Get from database
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 3. Cache the result
	if user != nil {
		go func() {
			u.cacheRepo.SetUserCache(id, user)
		}()
	}

	return user, nil
}

func (u *userUsecase) UpdateUser(user *entity.User) error {
	// 1. Update in database
	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	// 2. Update cache
	go func() {
		u.cacheRepo.SetUserCache(uint(user.ID), user)
	}()

	// 3. Publish event
	go func() {
		u.messageRepo.PublishUserUpdated(user)
	}()

	return nil
}

func (u *userUsecase) DeleteUser(id uint) error {
	// 1. Delete from database
	if err := u.userRepo.Delete(id); err != nil {
		return err
	}

	// 2. Delete from cache
	go func() {
		u.cacheRepo.DeleteUserCache(id)
	}()

	// 3. Publish event
	go func() {
		u.messageRepo.PublishUserDeleted(id)
	}()

	return nil
}

func (u *userUsecase) GetUserByEmail(email string) (*entity.User, error) {
	return u.userRepo.FindByEmail(email)
}

func (u *userUsecase) ValidateUserCredentials(email, password string) (*entity.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return user, nil
}