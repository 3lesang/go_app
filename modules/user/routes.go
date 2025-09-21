package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router) {
	users := r.Group("/users")
	users.Get("/", GetUsersHandler)
}
