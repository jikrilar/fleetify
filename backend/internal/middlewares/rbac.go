package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/responses"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := CurrentUser(c)
		if !ok {
			return responses.Fail(c, fiber.StatusUnauthorized, "User tidak ditemukan", "UNAUTHORIZED")
		}

		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return responses.Fail(c, fiber.StatusForbidden, "Role tidak memiliki akses ke fitur ini", "FORBIDDEN")
	}
}
