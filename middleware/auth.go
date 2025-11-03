package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewJWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.NoContent(http.StatusUnauthorized)
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.NoContent(http.StatusUnauthorized)
			}
			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}

			var userID uint
			if uid, ok := claims["user_id"].(float64); ok {
				userID = uint(uid)
			} else if uid, ok := claims["user_id"].(int); ok {
				userID = uint(uid)
			} else {

				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "user_id not found in token",
				})
			}
			c.Set("user_id", userID)

			if typeUser, ok := claims["type_user"].(string); ok {
				c.Set("type_user", typeUser)
			}
			return next(c)
		}
	}
}
