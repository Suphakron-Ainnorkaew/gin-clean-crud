package delivery

import "github.com/labstack/echo/v4"

func RegisterPublicUserRoutes(g *echo.Group, h *Handler) {
	g.POST("/users", h.CreateUser)
	g.POST("/login", h.Login)
}
