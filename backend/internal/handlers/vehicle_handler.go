package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/repositories"
	"github.com/jikrilar/fleetify/backend/internal/responses"
)

type VehicleHandler struct {
	vehicles *repositories.VehicleRepository
}

func NewVehicleHandler(vehicles *repositories.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{vehicles: vehicles}
}

func (h *VehicleHandler) List(c fiber.Ctx) error {
	vehicles, err := h.vehicles.FindAll()
	if err != nil {
		return responses.Fail(c, fiber.StatusInternalServerError, "Gagal mengambil kendaraan", "SERVER_ERROR")
	}
	return responses.Success(c, fiber.StatusOK, "Kendaraan berhasil diambil", vehicles)
}
