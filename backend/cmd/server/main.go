package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
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
	app.Get("/*", static.New(cfg.FrontendDir))

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
