package delivery

import (
	"go-clean-api/entity"
	"go-clean-api/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterAuthUserRoutes(g *echo.Group, h *Handler) {
	g.GET("/users", h.GetAllUsers, middleware.RequireRoleFromJWT(entity.UserTypeAdmin))
	g.GET("/users/:id", h.GetUserByID, middleware.RequireSelfOrAdmin())
	g.PUT("/users/:id", h.UpdateUser, middleware.RequireRoleFromJWT(entity.UserTypeAdmin))
	g.POST("/profile", h.ProfileUser)
	g.DELETE("/users/:id", h.DeleteUser)
}