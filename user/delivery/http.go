// user/delivery/http.go
package delivery

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase   domain.UserUsecase
	jwtSecret string
}

func NewHandler(e *echo.Group, usecase domain.UserUsecase, jwtSecret string) *Handler {
	handler := &Handler{
		usecase:   usecase,
		jwtSecret: jwtSecret,
	}

	e.POST("/users", handler.CreateUser)
	e.GET("/users", handler.GetAllUsers)
	e.GET("/users/:id", handler.GetUserByID)
	e.PUT("/users/:id", handler.UpdateUser)
	e.DELETE("/users/:id", handler.DeleteUser)

	e.POST("/login", handler.Login)
	e.POST("/profile", handler.ProfileUser)

	return handler
}

func (h *Handler) parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// POST /users
func (h *Handler) CreateUser(c echo.Context) error {
	var user entity.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password are required"})
	}

	if err := h.usecase.CreateUser(&user); err != nil {
		if err.Error() == "email already in use" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

// GET /users
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := h.usecase.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// GET /users/:id
func (h *Handler) GetUserByID(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// PUT /users/:id
func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.usecase.UpdateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// DELETE /users/:id
func (h *Handler) DeleteUser(c echo.Context) error {
	id, err := h.parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.usecase.DeleteUser(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
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

	secret := h.jwtSecret
	if secret == "" {
		secret = "secret"
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}

	sub, ok := claims["sub"]
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing subject claim"})
	}

	var userID uint
	switch v := sub.(type) {
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	user.Password = ""
	return c.JSON(http.StatusOK, user)
}
