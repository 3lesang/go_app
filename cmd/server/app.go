package server

import (
	"log"

	_ "app/docs"
	"app/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "github.com/mattn/go-sqlite3"
)

func Serve() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(logger.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the API 🚀",
		})
	})

	router.Setup(app)
	log.Println("Server started on port 8080")
	log.Fatal(app.Listen(":8080"))
}
