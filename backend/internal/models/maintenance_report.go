package models

import "time"

const (
	StatusPendingApproval = "PENDING_APPROVAL"
	StatusApproved        = "APPROVED"
	StatusCompleted       = "COMPLETED"
)

type MaintenanceReport struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	VehicleID    uint         `json:"vehicle_id"`
	CreatedBy    uint         `json:"created_by"`
	Odometer     uint         `json:"odometer"`
	Complaint    string       `json:"complaint"`
	Status       string       `json:"status" gorm:"type:enum('PENDING_APPROVAL','APPROVED','COMPLETED')"`
	InitialPhoto string       `json:"initial_photo"`
	ProofPhoto   string       `json:"proof_photo"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	User         User         `json:"user" gorm:"foreignKey:CreatedBy"`
	Vehicle      Vehicle      `json:"vehicle" gorm:"foreignKey:VehicleID"`
	Items        []ReportItem `json:"items" gorm:"foreignKey:ReportID"`
}
