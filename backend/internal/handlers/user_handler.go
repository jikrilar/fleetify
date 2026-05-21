package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/repositories"
	"github.com/jikrilar/fleetify/backend/internal/responses"
)

type UserHandler struct {
	users *repositories.UserRepository
}

func NewUserHandler(users *repositories.UserRepository) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) List(c fiber.Ctx) error {
	users, err := h.users.FindAll()
	if err != nil {
		return responses.Fail(c, fiber.StatusInternalServerError, "Gagal mengambil user", "SERVER_ERROR")
	}
	return responses.Success(c, fiber.StatusOK, "User berhasil diambil", users)
}
