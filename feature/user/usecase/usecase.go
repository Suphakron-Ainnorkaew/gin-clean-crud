package usecase

import (
	"errors"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewUserUsecase(repo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	return &userUsecase{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) CreateUser(user *entity.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Email == "" || user.Password == "" {
		return errors.New("email and password are required")
	}

	existing, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	if err := u.userRepo.Create(user); err != nil {
		return err
	}
	return nil

}

func (u *userUsecase) GetAllUsers() ([]entity.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) GetUserByID(id uint) (*entity.User, error) {

	if u == nil {
		return nil, errors.New("user usecase is nil")
	}
	if u.userRepo == nil {
		return nil, errors.New("user repository is not initialized")
	}

	return u.userRepo.FindByID(id)

}

func (u *userUsecase) UpdateUser(user *entity.User) error {
	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) DeleteUser(id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		return err
	}

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

func (u *userUsecase) Login(email, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"type_user": user.TypeUser,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
