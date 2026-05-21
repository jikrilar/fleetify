package repositories

import (
	"github.com/jikrilar/fleetify/backend/internal/models"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) DB() *gorm.DB {
	return r.db
}

func (r *ReportRepository) FindAll(status string) ([]models.MaintenanceReport, error) {
	var reports []models.MaintenanceReport
	query := r.db.Preload("User").Preload("Vehicle").Preload("Items.MasterItem").Order("created_at DESC")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&reports).Error
	return reports, err
}

func (r *ReportRepository) FindByID(id uint) (*models.MaintenanceReport, error) {
	var report models.MaintenanceReport
	err := r.db.Preload("User").Preload("Vehicle").Preload("Items.MasterItem").First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}
