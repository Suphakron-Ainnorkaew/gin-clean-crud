package delivery

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"go-clean-api/middleware"
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

func NewHandler(e *echo.Group, usecase domain.UserUsecase, cfg config.ToolsConfig) *Handler {

	e.Use(middleware.LoggingMiddleware(cfg.Logrus))

	handler := &Handler{
		usecase: usecase,
		cfg:     cfg,
	}

	return handler
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

	log := c.Get("logger").(*logrus.Entry)

	users, err := h.usecase.GetAllUsers(log)
	if err != nil {
		log.WithError(err).Error("Failed to get all user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// GET /users/:id
func (h *Handler) GetUserByID(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	user, err := h.usecase.GetUserByID(log, id)
	if err != nil {
		log.WithError(err).Error("Failed to Get user id in GetUserByID")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// PUT /users/:id
func (h *Handler) UpdateUser(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	user, err := h.usecase.GetUserByID(log, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.usecase.UpdateUser(log, user); err != nil {
		log.WithError(err).Error("Failed to update user in UpdateUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// DELETE /users/:id
func (h *Handler) DeleteUser(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	if err := h.usecase.DeleteUser(log, id); err != nil {
		log.WithError(err).Error("Failed to delete shop in DeleteUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GET Profile
func (h *Handler) ProfileUser(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

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

	user, err := h.usecase.GetUserByID(log, userID)
	if err != nil {
		log.WithError(err).Error("Failed to load profile in ProfileUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// POST /users
func (h *Handler) CreateUser(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	var user entity.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password are required"})
	}

	if err := h.usecase.CreateUser(log, &user); err != nil {
		if err.Error() == "email already in use" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		log.WithError(err).Error("Failed to create user in CreateUser")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {

	log := c.Get("logger").(*logrus.Entry)

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	token, err := h.usecase.Login(log, req.Email, req.Password)
	if err != nil {
		log.WithError(err).Error("Failed to login in Login")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
