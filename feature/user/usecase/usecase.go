package usecase

import (
	"errors"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
	cfg      config.ToolsConfig
}

func NewUserUsecase(repo domain.UserRepository, cfg config.ToolsConfig) domain.UserUsecase {
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
		u.cfg.Logrus.WithError(err).WithField("email", user.Email).Error("failed to check existing email")
		return err
	}
	if existing != nil {
		u.cfg.Logrus.WithField("email", user.Email).Info("email already in use")
		return errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	if err := u.userRepo.Create(user); err != nil {
		u.cfg.Logrus.WithError(err).Error("failed to create user")
		return err
	}
	u.cfg.Logrus.WithField("userID", user.ID).Info("user created")
	return nil

}

func (u *userUsecase) GetAllUsers() ([]entity.User, error) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		u.cfg.Logrus.WithError(err).Error("failed to get all users")
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) GetUserByID(id uint) (*entity.User, error) {

	if u == nil {
		return nil, errors.New("user usecase is nil")
	}
	if u.userRepo == nil {
		return nil, errors.New("user repository is not initialized")
	}

	user, err := u.userRepo.FindByID(id)
	if err != nil {
		u.cfg.Logrus.WithError(err).WithField("userID", id).Error("failed to get user by id")
		return nil, err
	}
	return user, nil

}

func (u *userUsecase) UpdateUser(user *entity.User) error {
	if err := u.userRepo.Update(user); err != nil {
		u.cfg.Logrus.WithError(err).WithField("userID", user.ID).Error("failed to update user")
		return err
	}
	u.cfg.Logrus.WithField("userID", user.ID).Info("user updated")
	return nil
}

func (u *userUsecase) DeleteUser(id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		u.cfg.Logrus.WithError(err).WithField("userID", id).Error("failed to delete user")
		return err
	}
	u.cfg.Logrus.WithField("userID", id).Info("user deleted")
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
		u.cfg.Logrus.WithField("email", email).Warn("login failed: user not found")
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		u.cfg.Logrus.WithField("email", email).Warn("login failed: wrong password")
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
		u.cfg.Logrus.WithError(err).Error("failed to sign token")
		return "", err
	}
	return signed, nil
}
