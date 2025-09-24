package router

import (
	"app/internal/user"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	user.RegisterRoutes(api)
}
