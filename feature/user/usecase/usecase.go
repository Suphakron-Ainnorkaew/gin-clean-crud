package usecase

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
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

func (u *userUsecase) CreateUser(user *entity.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if user.Email == "" || user.Password == "" {
		return errors.New("email and password are required")
	}

	existing, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		err := errors.New("[Usecase.CreateUser]: failed to find email user")
		log.Warn(err)
		return err
	}
	if existing != nil {
		return errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err := errors.New("[Usecase.CreateUser]: failed to generate password user")
		log.Warn(err)
		return err
	}
	user.Password = string(hashed)

	if err := u.userRepo.Create(user); err != nil {
		err := errors.Wrap(err, "[Usecase.CreateUser]: failed to create user")
		log.Error(err)
		return err
	}

	return nil

}

func (u *userUsecase) GetAllUsers() ([]entity.User, error) {
	user, err := u.userRepo.FindAll()
	if err != nil {
		wrappedErr := errors.Wrap(err, "[Usecase.GetAllUsers]: failed to get all user")
		log.Error(wrappedErr)
		return nil, wrappedErr
	}
	return user, nil
}

func (u *userUsecase) GetUserByID(id uint) (*entity.User, error) {
	user, err := u.userRepo.FindByID(id)

	if u == nil {
		return nil, errors.New("user usecase is nil")
	}
	if u.userRepo == nil {
		return nil, errors.New("user repository is not initialized")
	}

	if err != nil {
		err := errors.Wrap(err, "[Usecase.GetUserByID]: failed to get user id")
		log.Error(err)
		return nil, err
	}

	return user, nil

}

func (u *userUsecase) UpdateUser(user *entity.User) error {

	if err := u.userRepo.Update(user); err != nil {
		err := errors.Wrap(err, "[Usecase.UpdateUser]: failed to update user")
		log.Error(err)
		return err
	}

	return nil
}

func (u *userUsecase) DeleteUser(id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		err := errors.Wrap(err, "[Usecase.DeleteUser]: failed to delete user")
		log.Error(err)
		return err
	}

	return nil
}

func (u *userUsecase) GetUserByEmail(email string) (*entity.User, error) {

	user, err := u.userRepo.FindByEmail(email)

	if err != nil {
		err := errors.Wrap(err, "[Usecase.GetUserByEmail]: failed to get email user")
		log.Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) ValidateUserCredentials(email, password string) (*entity.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		err := errors.Wrap(err, "[Usecase.ValidateUserCredentials]: failed to validate user")
		log.Error(err)
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
		err := errors.Wrap(err, "[Usecase.Login]: failed to find email user")
		log.Warn(err)
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
		err := errors.Wrap(err, "[Usecase.Login]: failed to login user")
		log.Error(err)
		return "", err
	}
	return signed, nil
}
