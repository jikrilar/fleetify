package middlewares

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/models"
	"github.com/jikrilar/fleetify/backend/internal/repositories"
	"github.com/jikrilar/fleetify/backend/internal/responses"
	"gorm.io/gorm"
)

const userContextKey = "current_user"

func UserContext(users *repositories.UserRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		headerValue := c.Get("X-User-ID")
		if headerValue == "" {
			return responses.Fail(c, fiber.StatusUnauthorized, "Header X-User-ID wajib diisi", "UNAUTHORIZED")
		}

		parsed, err := strconv.ParseUint(headerValue, 10, 64)
		if err != nil {
			return responses.Fail(c, fiber.StatusUnauthorized, "Header X-User-ID harus berupa angka", "UNAUTHORIZED")
		}

		user, err := users.FindByID(uint(parsed))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return responses.Fail(c, fiber.StatusUnauthorized, "User tidak ditemukan", "UNAUTHORIZED")
			}
			return responses.Fail(c, fiber.StatusInternalServerError, "Gagal memuat user", "SERVER_ERROR")
		}

		c.Locals(userContextKey, *user)
		return c.Next()
	}
}

func CurrentUser(c fiber.Ctx) (models.User, bool) {
	user, ok := c.Locals(userContextKey).(models.User)
	return user, ok
}
