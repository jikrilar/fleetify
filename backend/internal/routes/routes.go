package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jikrilar/fleetify/backend/internal/handlers"
	"github.com/jikrilar/fleetify/backend/internal/middlewares"
	"github.com/jikrilar/fleetify/backend/internal/models"
	"github.com/jikrilar/fleetify/backend/internal/repositories"
	"github.com/jikrilar/fleetify/backend/internal/services"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB, webhookURL string) {
	userRepo := repositories.NewUserRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	userHandler := handlers.NewUserHandler(userRepo)
	vehicleHandler := handlers.NewVehicleHandler(vehicleRepo)
	itemHandler := handlers.NewItemHandler(itemRepo)
	reportHandler := handlers.NewReportHandler(services.NewReportService(reportRepo, services.NewWebhookService(webhookURL)))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	api.Get("/users", userHandler.List)

	protected := api.Group("/", middlewares.UserContext(userRepo))
	protected.Get("/vehicles", middlewares.RequireRole(models.RoleSA, models.RoleApproval), vehicleHandler.List)
	protected.Get("/master-items", middlewares.RequireRole(models.RoleSA, models.RoleApproval), itemHandler.List)
	protected.Get("/reports", middlewares.RequireRole(models.RoleSA, models.RoleApproval), reportHandler.List)
	protected.Get("/reports/:id", middlewares.RequireRole(models.RoleSA, models.RoleApproval), reportHandler.Detail)
	protected.Post("/reports", middlewares.RequireRole(models.RoleSA), reportHandler.Create)
	protected.Patch("/reports/:id/approve", middlewares.RequireRole(models.RoleApproval), reportHandler.Approve)
	protected.Patch("/reports/:id/complete", middlewares.RequireRole(models.RoleSA), reportHandler.Complete)
}
