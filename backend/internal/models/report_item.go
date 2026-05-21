package models

import "time"

type ReportItem struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	ReportID      uint       `json:"report_id"`
	ItemID        uint       `json:"item_id"`
	Quantity      uint       `json:"quantity"`
	PriceSnapshot float64    `json:"price_snapshot"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	MasterItem    MasterItem `json:"master_item" gorm:"foreignKey:ItemID"`
}
