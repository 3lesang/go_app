package user

import (
	"app/internal/db"
	auth_db "app/internal/db/auth"
	"context"
	"crypto/md5"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// GetUsersHandler godoc
// @Summary      Get user list
// @Description  Returns a list of users
// @Tags         users
// @Security BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /users [get]
func GetUsersHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	users, err := db.AuthQueries.ListUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": users,
	})
}

// GetUserHandler godoc
// @Summary      Get a user
// @Description  Returns a user by ID
// @Tags         users
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
func GetUserHandler(c *fiber.Ctx) error {
	param := c.Params("id")

	ctx := context.Background()
	user, err := db.AuthQueries.GetUser(ctx, pgtype.UUID{Bytes: md5.Sum([]byte(param)), Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// CreateUserHandler godoc
// @Summary      Create a new user
// @Description  Creates a new user and returns the created user
// @Tags         users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users [post]
func CreateUserHandler(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.AuthQueries.CreateUser(ctx, auth_db.CreateUserParams{
		Name:     req.Name,
		Phone:    pgtype.Text{String: req.Phone, Valid: true},
		Email:    pgtype.Text{String: req.Email, Valid: true},
		Username: pgtype.Text{String: req.Username, Valid: true},
		Password: string(hash),
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
		"user": fiber.Map{
			"name":  req.Name,
			"email": req.Email,
		},
	})
}

// DeleteUsersHandler godoc
// @Summary      Delete multiple users
// @Description  Deletes multiple users by their IDs
// @Tags         users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteUsersRequest  true  "List of user IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users [delete]
func DeleteUsersHandler(c *fiber.Ctx) error {
	var req DeleteUsersRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no user IDs provided",
		})
	}

	ctx := context.Background()
	if err := db.AuthQueries.DeleteUsers(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "users deleted successfully",
		"count":   len(req.IDs),
	})
}
