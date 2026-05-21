package repositories

import (
	"github.com/jikrilar/fleetify/backend/internal/models"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) FindAll() ([]models.MasterItem, error) {
	var items []models.MasterItem
	err := r.db.Order("item_name ASC").Find(&items).Error
	return items, err
}
