package delivery

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase domain.UserUsecase
	cfg     config.ToolsConfig
}

func NewHandler(usecase domain.UserUsecase, cfg config.ToolsConfig) *Handler {
	return &Handler{
		usecase: usecase,
		cfg:     cfg,
	}
}

func (h *Handler) parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GET /users
func (h *Handler) GetAllUsers(c echo.Context) error {
	h.cfg.Logrus.Info("getall user request")
	users, err := h.usecase.GetAllUsers()
	if err != nil {
		h.cfg.Logrus.WithError(err).Error("failed to get all users in GetAllUsers")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.cfg.Logrus.WithField("count", len(users)).Info("getall-user success")
	return c.JSON(http.StatusOK, users)
}

// GET /users/:id
func (h *Handler) GetUserByID(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid user ID in GetUserByID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		h.cfg.Logrus.WithError(err).Error("GetUserByID Failed")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		h.cfg.Logrus.WithField("userID", id).Warn("user not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	h.cfg.Logrus.WithField("userID", id).Info("get-userbyid success")
	return c.JSON(http.StatusOK, user)
}

// PUT /users/:id
func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid user id in UpdateUser")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		h.cfg.Logrus.WithError(err).Error("failed to get user by id in UpdateUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.usecase.UpdateUser(user); err != nil {
		h.cfg.Logrus.WithError(err).Error("failed to update user in UpdateUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.cfg.Logrus.WithField("userID", id).Info("update-user success")
	return c.JSON(http.StatusOK, user)
}

// DELETE /users/:id
func (h *Handler) DeleteUser(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid user id in DeleteUser")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	if err := h.usecase.DeleteUser(id); err != nil {
		h.cfg.Logrus.WithError(err).Error("Failed to delete user for DeleteUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.cfg.Logrus.WithField("userID", id).Info("delete-user success")
	return c.NoContent(http.StatusNoContent)
}

// GET Profile
func (h *Handler) ProfileUser(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
	}
	tokenString := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing oken"})
	}

	secret := h.cfg.JWTSecret
	if secret == "" {
		secret = "secret"
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		h.cfg.Logrus.WithError(err).Warn("invalid JWT in ProfileUser")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}

	uidClaim, ok := claims["user_id"]
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing subject claim"})
	}

	var userID uint
	switch v := uidClaim.(type) {
	case float64:
		userID = uint(v)
	case string:
		id64, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid subject claim"})
		}
		userID = uint(id64)
	default:
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid subject claim type"})
	}

	user, err := h.usecase.GetUserByID(userID)
	if err != nil {
		h.cfg.Logrus.WithError(err).Error("failed to get user by id in ProfileUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		h.cfg.Logrus.WithField("userID", userID).Warn("user not found in ProfileUser")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// POST /users
func (h *Handler) CreateUser(c echo.Context) error {
	var user entity.User

	log := h.cfg.Logrus.WithFields(logrus.Fields{
		"endpoint": "POST /users",
		"method":   c.Request().Method,
		"path":     c.Path(),
	})

	log.Info("Request started")

	if err := c.Bind(&user); err != nil {
		log.WithError(err).Warn("Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if user.Email == "" || user.Password == "" {
		log.Warn("Missing required fields: email or password")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password are required"})
	}

	if err := h.usecase.CreateUser(&user); err != nil {
		if err.Error() == "email already in use" {
			log.WithField("email", user.Email).Info("Email already in use")
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}

		log.WithError(err).Error("Failed to create user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	user.Password = ""
	log.WithField("email", user.Email).Info("User created successfully")

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		h.cfg.Logrus.WithError(err).Warn("invalid user in Login")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		h.cfg.Logrus.WithError(err).WithField("email", req.Email).Warn("login failed")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
