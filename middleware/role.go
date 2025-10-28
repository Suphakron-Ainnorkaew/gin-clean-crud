package middleware

import (
	"net/http"

	"go-clean-api/entity"

	"github.com/labstack/echo/v4"
)

type UserFetcher func(id uint) (*entity.User, error)

func RequireRole(fetch UserFetcher, allowedRoles ...entity.UserType) echo.MiddlewareFunc {
	roleSet := map[entity.UserType]struct{}{}
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			v := c.Get("user_id")
			if v == nil {
				return c.NoContent(http.StatusUnauthorized)
			}
			uid, ok := v.(uint)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}
			user, err := fetch(uid)
			if err != nil || user == nil {
				return c.NoContent(http.StatusUnauthorized)
			}
			uRole := entity.UserType(user.TypeUser)
			if _, ok := roleSet[uRole]; !ok {
				return c.NoContent(http.StatusForbidden)
			}
			return next(c)
		}
	}
}

// RequireRoleFromJWT checks type_user from JWT claims (faster, no DB query)
// Use this when you trust the JWT token and don't need to verify against DB
func RequireRoleFromJWT(allowedRoles ...entity.UserType) echo.MiddlewareFunc {
	roleSet := map[entity.UserType]struct{}{}
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			typeUserVal := c.Get("type_user")
			if typeUserVal == nil {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "type_user not found in token",
				})
			}
			typeUser, ok := typeUserVal.(string)
			if !ok {
				return c.NoContent(http.StatusForbidden)
			}
			uRole := entity.UserType(typeUser)
			if _, ok := roleSet[uRole]; !ok {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "insufficient permissions",
					"required_role": "shop",
					"your_role": typeUser,
				})
			}
			return next(c)
		}
	}
}
