package user

import (
	"app/internal/db"

	"context"

	"github.com/gofiber/fiber/v2"
)

// GetUsersHandler godoc
// @Summary      Get user list
// @Description  Returns a list of users
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /users [get]
func GetUsersHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	users, err := db.Queries.ListUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	data := fiber.Map{
		"data": users,
	}

	return c.JSON(data)
}
