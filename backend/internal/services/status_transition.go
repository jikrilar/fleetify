package services

import "github.com/jikrilar/fleetify/backend/internal/models"

func canApprove(status string) bool {
	return status == models.StatusPendingApproval
}

func canComplete(status string) bool {
	return status == models.StatusApproved
}
