package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/config"
	"github.com/jikrilar/fleetify/backend/internal/database"
	"github.com/jikrilar/fleetify/backend/internal/routes"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}

	app := fiber.New()
	routes.Register(app, db, cfg.WebhookURL)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
