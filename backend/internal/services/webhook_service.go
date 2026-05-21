package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jikrilar/fleetify/backend/internal/models"
)

type WebhookService struct {
	url string
}

type webhookPayload struct {
	Event               string    `json:"event"`
	ReportID            uint      `json:"report_id"`
	Status              string    `json:"status"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func NewWebhookService(url string) *WebhookService {
	return &WebhookService{url: url}
}

func (s *WebhookService) SendAsync(event string, report models.MaintenanceReport) {
	if s.url == "" {
		return
	}

	go func() {
		payload := webhookPayload{
			Event:               event,
			ReportID:            report.ID,
			Status:              report.Status,
			VehicleLicensePlate: report.Vehicle.LicensePlate,
			UpdatedAt:           report.UpdatedAt,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("gagal membuat payload webhook: %v", err)
			return
		}

		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Post(s.url, "application/json", bytes.NewReader(body))
		if err != nil {
			log.Printf("webhook gagal dikirim: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			log.Printf("webhook mendapat status tidak sukses: %d", resp.StatusCode)
		}
	}()
}
