// user/delivery/http.go
package delivery

import (
	"go-clean-api/domain"
	"go-clean-api/entity"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.UserUsecase
}

func NewHandler(e *echo.Group, usecase domain.UserUsecase) *Handler {
	handler := &Handler{
		usecase: usecase,
	}

	e.POST("/users", handler.CreateUser)
	e.GET("/users", handler.GetAllUsers)
	e.GET("/users/:id", handler.GetUserByID)
	e.PUT("/users/:id", handler.UpdateUser)
	e.DELETE("/users/:id", handler.DeleteUser)

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

	if err := h.usecase.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

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