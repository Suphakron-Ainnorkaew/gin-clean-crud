package usecase

import (
	"errors"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
	cfg      config.ToolsConfig
}

func NewUserUsecase(repo domain.UserRepository,
	cfg config.ToolsConfig,
) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
		cfg:      cfg,
	}
}

func (u *userUsecase) CreateUser(log *logrus.Entry, user *entity.User) error {
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
		log.WithError(err).Error("Failed to create shop in repo")
		return err
	}
	return nil

}

func (u *userUsecase) GetAllUsers(log *logrus.Entry) ([]entity.User, error) {
	user, err := u.userRepo.FindAll()
	if err != nil {
		log.WithError(err).Error("Failed to get-all user from repo")
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetUserByID(log *logrus.Entry, id uint) (*entity.User, error) {

	log = log.WithField("user_id", id)
	user, err := u.userRepo.FindByID(id)

	if u == nil {
		return nil, errors.New("user usecase is nil")
	}
	if u.userRepo == nil {
		return nil, errors.New("user repository is not initialized")
	}

	if err != nil {
		log.WithError(err).Error("Failed to get user id from repo")
		return nil, err
	}

	return user, nil

}

func (u *userUsecase) UpdateUser(log *logrus.Entry, user *entity.User) error {

	log = log.WithField("user_id", user.ID)

	if err := u.userRepo.Update(user); err != nil {
		log.WithError(err).Error("Failed to update user in repo")
		return err
	}

	return nil
}

func (u *userUsecase) DeleteUser(log *logrus.Entry, id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		log.WithError(err).Error("Failed to delete user in repo")
		return err
	}

	return nil
}

func (u *userUsecase) GetUserByEmail(log *logrus.Entry, email string) (*entity.User, error) {

	user, err := u.userRepo.FindByEmail(email)

	if err != nil {
		log.WithError(err).Error("Failed to get email user in repo")
		return user, nil
	}
	return user, nil
}

func (u *userUsecase) ValidateUserCredentials(log *logrus.Entry, email, password string) (*entity.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		log.WithError(err).Error("Failed to validate user from repo")
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (u *userUsecase) Login(log *logrus.Entry, email, password string) (string, error) {
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
	signed, err := token.SignedString([]byte(u.cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
