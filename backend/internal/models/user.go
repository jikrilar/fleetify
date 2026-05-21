package models

import "time"

const (
	RoleSA       = "SA"
	RoleApproval = "APPROVAL"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Role      string    `json:"role" gorm:"type:enum('SA','APPROVAL')"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
