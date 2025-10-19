package server

import (
	"log"

	_ "app/docs"
	"app/internal/router"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func Serve() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	app.Use(logger.New())
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:                  "/swagger/doc.json",
		PersistAuthorization: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the API 🚀",
		})
	})

	router.Init(app)
	log.Println("Server started on port 8080")
	log.Fatal(app.Listen(":8080"))
}
