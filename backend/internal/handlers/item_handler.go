package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/repositories"
	"github.com/jikrilar/fleetify/backend/internal/responses"
)

type ItemHandler struct {
	items *repositories.ItemRepository
}

func NewItemHandler(items *repositories.ItemRepository) *ItemHandler {
	return &ItemHandler{items: items}
}

func (h *ItemHandler) List(c fiber.Ctx) error {
	items, err := h.items.FindAll()
	if err != nil {
		return responses.Fail(c, fiber.StatusInternalServerError, "Gagal mengambil item master", "SERVER_ERROR")
	}
	return responses.Success(c, fiber.StatusOK, "Item master berhasil diambil", items)
}
