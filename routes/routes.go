package routes

import (
	"github.com/gofiber/fiber/v2"

	"app/modules/user"
)

func Register(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	user.RegisterRoutes(v1)
}
