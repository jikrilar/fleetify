package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/middlewares"
	"github.com/jikrilar/fleetify/backend/internal/responses"
	"github.com/jikrilar/fleetify/backend/internal/services"
	"gorm.io/gorm"
)

type ReportHandler struct {
	reports *services.ReportService
}

func NewReportHandler(reports *services.ReportService) *ReportHandler {
	return &ReportHandler{reports: reports}
}

func (h *ReportHandler) List(c fiber.Ctx) error {
	status := c.Query("status")
	reports, err := h.reports.ListReports(status)
	if err != nil {
		return responses.Fail(c, fiber.StatusInternalServerError, "Gagal mengambil laporan", "SERVER_ERROR")
	}
	return responses.Success(c, fiber.StatusOK, "Laporan berhasil diambil", reports)
}

func (h *ReportHandler) Detail(c fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return responses.Fail(c, fiber.StatusBadRequest, "ID laporan tidak valid", "BAD_REQUEST")
	}

	report, err := h.reports.GetReport(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.Fail(c, fiber.StatusNotFound, "Laporan tidak ditemukan", "NOT_FOUND")
		}
		return responses.Fail(c, fiber.StatusInternalServerError, "Gagal mengambil detail laporan", "SERVER_ERROR")
	}
	return responses.Success(c, fiber.StatusOK, "Detail laporan berhasil diambil", report)
}

func (h *ReportHandler) Create(c fiber.Ctx) error {
	var req services.CreateReportRequest
	if err := c.Bind().Body(&req); err != nil {
		return responses.Fail(c, fiber.StatusBadRequest, "Request tidak valid", "BAD_REQUEST")
	}

	user, ok := middlewares.CurrentUser(c)
	if !ok {
		return responses.Fail(c, fiber.StatusUnauthorized, "User tidak ditemukan", "UNAUTHORIZED")
	}

	report, err := h.reports.CreateReport(user, req)
	if err != nil {
		return responses.Fail(c, fiber.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
	}
	return responses.Success(c, fiber.StatusCreated, "Laporan berhasil dibuat", report)
}

func (h *ReportHandler) Approve(c fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return responses.Fail(c, fiber.StatusBadRequest, "ID laporan tidak valid", "BAD_REQUEST")
	}

	report, err := h.reports.ApproveReport(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.Fail(c, fiber.StatusNotFound, "Laporan tidak ditemukan", "NOT_FOUND")
		}
		return responses.Fail(c, fiber.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
	}
	return responses.Success(c, fiber.StatusOK, "Laporan berhasil disetujui", report)
}

func (h *ReportHandler) Complete(c fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return responses.Fail(c, fiber.StatusBadRequest, "ID laporan tidak valid", "BAD_REQUEST")
	}

	var req services.CompleteReportRequest
	if err := c.Bind().Body(&req); err != nil {
		return responses.Fail(c, fiber.StatusBadRequest, "Request tidak valid", "BAD_REQUEST")
	}

	report, err := h.reports.CompleteReport(id, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.Fail(c, fiber.StatusNotFound, "Laporan tidak ditemukan", "NOT_FOUND")
		}
		return responses.Fail(c, fiber.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
	}
	return responses.Success(c, fiber.StatusOK, "Laporan berhasil diselesaikan", report)
}

func parseID(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
