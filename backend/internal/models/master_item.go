package models

import "time"

const (
	ItemTypePart    = "PART"
	ItemTypeService = "SERVICE"
)

type MasterItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ItemName  string    `json:"item_name"`
	Type      string    `json:"type" gorm:"type:enum('PART','SERVICE')"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
