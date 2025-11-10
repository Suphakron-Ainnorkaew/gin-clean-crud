package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(user *entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.Create]: unable to create user")
	}
	return nil
}

func (r *postgresUserRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "[UserRepository.FindAll]: unable to findall user")
	}
	return users, nil
}

func (r *postgresUserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User

	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[UserRepository.FindByID]: unable to find user by id")
	}

	return &user, nil
}

func (r *postgresUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[UserRepository.FindByEmail]: unable to find user by email")
	}
	return &user, nil
}

func (r *postgresUserRepository) Update(user *entity.User) error {
	if err := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.Update]: unable to update user")
	}
	return nil
}

func (r *postgresUserRepository) Delete(id uint) error {
	if err := r.db.Delete(&entity.User{}, id).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.Delete]: unable to delete user")
	}
	return nil
}
