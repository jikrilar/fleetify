package services

import (
	"testing"

	"github.com/jikrilar/fleetify/backend/internal/models"
)

func TestCanApprove(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{name: "pending approval bisa disetujui", status: models.StatusPendingApproval, want: true},
		{name: "approved tidak bisa disetujui ulang", status: models.StatusApproved, want: false},
		{name: "completed tidak bisa disetujui ulang", status: models.StatusCompleted, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canApprove(tt.status); got != tt.want {
				t.Fatalf("canApprove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanComplete(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{name: "approved bisa diselesaikan", status: models.StatusApproved, want: true},
		{name: "pending approval belum bisa diselesaikan", status: models.StatusPendingApproval, want: false},
		{name: "completed tidak bisa diselesaikan ulang", status: models.StatusCompleted, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canComplete(tt.status); got != tt.want {
				t.Fatalf("canComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}
