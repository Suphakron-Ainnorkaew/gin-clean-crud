package delivery

import (
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"go-clean-api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type Handler struct {
	usecase domain.UserUsecase
	cfg     config.ToolsConfig
}

func NewHandler(e *echo.Group, usecase domain.UserUsecase, cfg config.ToolsConfig) *Handler {

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

	users, err := h.usecase.GetAllUsers()
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetAllUsers]: failed to get user")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	return c.JSON(http.StatusOK, users)
}

// GET /users/:id
func (h *Handler) GetUserByID(c echo.Context) error {

	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetUserByID]: invalid user id")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		err = errors.Wrap(err, "[Handler.GetUserByID]: failed to get user id")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	if user == nil {
		err = errors.Wrap(err, "[Handler.GetUserByID]: user not found")
		log.Warn(err)
		return c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// PUT /users/:id
func (h *Handler) UpdateUser(c echo.Context) error {

	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: invalid user id")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: failed to get user id")
		log.Warn(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	if user == nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: user not found")
		log.Warn(err)
		return c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	if err := c.Bind(user); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: invalid request body")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	if err := h.usecase.UpdateUser(user); err != nil {
		err = errors.Wrap(err, "[Handler.UpdateUser]: failed to update user")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DELETE /users/:id
func (h *Handler) DeleteUser(c echo.Context) error {

	id, err := h.parseID(c)
	if err != nil {
		err = errors.Wrap(err, "[Handler.DeleteUser]: invalid user ID")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	if err := h.usecase.DeleteUser(id); err != nil {
		err = errors.Wrap(err, "[Handler.DeleteUser]: failed to delete user")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

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
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
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

	user, err := h.usecase.GetUserByID(userID)
	if err != nil {
		err = errors.Wrap(err, "[Handler.ProfileUser]: failed to show profile user")
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	if user == nil {
		err = errors.Wrap(err, "[Handler.ProfileUser]: user not found")
		log.Warn(err)
		return c.JSON(http.StatusNotFound, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// POST /users
func (h *Handler) CreateUser(c echo.Context) error {

	var user entity.User

	if err := c.Bind(&user); err != nil {
		err = errors.Wrap(err, "[Handler.CreateUser]: failed to bind request body")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	if user.Email == "" || user.Password == "" {
		err := errors.New("[Handler.CreateUser]: email or password missing")
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	if err := h.usecase.CreateUser(&user); err != nil {
		err = errors.Wrap(err, "[Handler.CreateUser]: failed to create user")
		log.Error(err)

		if err.Error() == "email already in use" {
			return c.JSON(http.StatusConflict, entity.ErrorResponse{
				Message: utils.StandardError(err),
			})
		}
		return c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		err = errors.Wrap(err, "[Handler.Login]: failed to bind request body")
		log.Warn(err)

		return c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		err = errors.Wrap(err, "[Handler.Login]: failed to login")
		log.Error(err)
		return c.JSON(http.StatusUnauthorized, entity.ErrorResponse{
			Message: utils.StandardError(err),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
