package user

import (
	"go-clean-api/internal/user/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindAll() ([]domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *gormUserRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Debug().Find(&users).Error
	return users, err
}

func (r *gormUserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *gormUserRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
