package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/jikrilar/fleetify/backend/internal/models"
	"github.com/jikrilar/fleetify/backend/internal/repositories"
	"gorm.io/gorm"
)

type CreateReportRequest struct {
	VehicleID    uint                    `json:"vehicle_id"`
	Odometer     uint                    `json:"odometer"`
	Complaint    string                  `json:"complaint"`
	InitialPhoto string                  `json:"initial_photo"`
	Items        []CreateReportItemInput `json:"items"`
}

type CreateReportItemInput struct {
	ItemID   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
}

type CompleteReportRequest struct {
	ProofPhoto string `json:"proof_photo"`
}

type ReportResponse struct {
	models.MaintenanceReport
	TotalEstimate float64 `json:"total_estimate"`
}

type ReportService struct {
	reports *repositories.ReportRepository
	webhook *WebhookService
}

func NewReportService(reports *repositories.ReportRepository, webhook *WebhookService) *ReportService {
	return &ReportService{reports: reports, webhook: webhook}
}

func (s *ReportService) ListReports(status string) ([]ReportResponse, error) {
	reports, err := s.reports.FindAll(status)
	if err != nil {
		return nil, err
	}
	return buildReportResponses(reports), nil
}

func (s *ReportService) GetReport(id uint) (*ReportResponse, error) {
	report, err := s.reports.FindByID(id)
	if err != nil {
		return nil, err
	}
	response := buildReportResponse(*report)
	return &response, nil
}

func (s *ReportService) CreateReport(user models.User, req CreateReportRequest) (*ReportResponse, error) {
	if req.VehicleID == 0 {
		return nil, errors.New("kendaraan wajib dipilih")
	}
	if req.Odometer == 0 {
		return nil, errors.New("odometer wajib lebih dari 0")
	}
	if req.Complaint == "" {
		return nil, errors.New("keluhan wajib diisi")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("minimal pilih satu item estimasi")
	}

	db := s.reports.DB()
	var createdID uint

	err := db.Transaction(func(tx *gorm.DB) error {
		var vehicle models.Vehicle
		if err := tx.First(&vehicle, req.VehicleID).Error; err != nil {
			return errors.New("kendaraan tidak ditemukan")
		}

		report := models.MaintenanceReport{
			VehicleID:    req.VehicleID,
			CreatedBy:    user.ID,
			Odometer:     req.Odometer,
			Complaint:    req.Complaint,
			Status:       models.StatusPendingApproval,
			InitialPhoto: req.InitialPhoto,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := tx.Create(&report).Error; err != nil {
			return err
		}

		for _, itemInput := range req.Items {
			if itemInput.ItemID == 0 || itemInput.Quantity == 0 {
				return errors.New("item dan quantity wajib valid")
			}

			var item models.MasterItem
			if err := tx.First(&item, itemInput.ItemID).Error; err != nil {
				return fmt.Errorf("item dengan id %d tidak ditemukan", itemInput.ItemID)
			}

			reportItem := models.ReportItem{
				ReportID:      report.ID,
				ItemID:        item.ID,
				Quantity:      itemInput.Quantity,
				PriceSnapshot: item.Price,
			}
			if err := tx.Create(&reportItem).Error; err != nil {
				return err
			}
		}

		createdID = report.ID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetReport(createdID)
}

func (s *ReportService) ApproveReport(id uint) (*ReportResponse, error) {
	db := s.reports.DB()
	var report models.MaintenanceReport

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&report, id).Error; err != nil {
			return err
		}
		if report.Status != models.StatusPendingApproval {
			return errors.New("hanya laporan PENDING_APPROVAL yang bisa disetujui")
		}
		report.Status = models.StatusApproved
		report.UpdatedAt = time.Now()
		return tx.Save(&report).Error
	})
	if err != nil {
		return nil, err
	}

	response, err := s.GetReport(id)
	if err != nil {
		return nil, err
	}
	s.webhook.SendAsync("REPORT_APPROVED", response.MaintenanceReport)
	return response, nil
}

func (s *ReportService) CompleteReport(id uint, req CompleteReportRequest) (*ReportResponse, error) {
	if req.ProofPhoto == "" {
		return nil, errors.New("foto bukti wajib diisi")
	}

	db := s.reports.DB()
	var report models.MaintenanceReport

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&report, id).Error; err != nil {
			return err
		}
		if report.Status != models.StatusApproved {
			return errors.New("hanya laporan APPROVED yang bisa diselesaikan")
		}
		report.Status = models.StatusCompleted
		report.ProofPhoto = req.ProofPhoto
		report.UpdatedAt = time.Now()
		return tx.Save(&report).Error
	})
	if err != nil {
		return nil, err
	}

	response, err := s.GetReport(id)
	if err != nil {
		return nil, err
	}
	s.webhook.SendAsync("REPORT_COMPLETED", response.MaintenanceReport)
	return response, nil
}

func buildReportResponses(reports []models.MaintenanceReport) []ReportResponse {
	responses := make([]ReportResponse, 0, len(reports))
	for _, report := range reports {
		responses = append(responses, buildReportResponse(report))
	}
	return responses
}

func buildReportResponse(report models.MaintenanceReport) ReportResponse {
	total := 0.0
	for _, item := range report.Items {
		total += item.PriceSnapshot * float64(item.Quantity)
	}
	return ReportResponse{MaintenanceReport: report, TotalEstimate: total}
}
