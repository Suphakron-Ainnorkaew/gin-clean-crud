// user/repository/postgres.go
package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *postgresUserRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *postgresUserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *postgresUserRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}