package responses

import "github.com/gofiber/fiber/v3"

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Success(c fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Fail(c fiber.Ctx, status int, message string, code string) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Message: message,
		Error:   code,
	})
}
