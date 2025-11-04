package server

import (
	"log"

	_ "app/docs"
	"app/internal/router"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost:3000,https://admin.senhome.vn,https://web-dev.senhome.vn,https://senhome.vn",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length, Authorization",
		AllowCredentials: true,
	}))
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
