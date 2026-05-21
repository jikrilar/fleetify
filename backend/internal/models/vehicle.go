package models

import "time"

type Vehicle struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	LicensePlate string    `json:"license_plate"`
	Model        string    `json:"model"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
