package repositories

import (
	"github.com/jikrilar/fleetify/backend/internal/models"
	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) FindAll() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := r.db.Order("license_plate ASC").Find(&vehicles).Error
	return vehicles, err
}
