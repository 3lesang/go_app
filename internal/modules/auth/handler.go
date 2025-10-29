package auth

import (
	"app/internal/db"
	auth_db "app/internal/db/auth"
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler godoc
// @Summary      User login
// @Description  Authenticates a user using email or username and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      LoginRequest  true  "Login request"
// @Success      200  {object}  map[string]interface{}  "Login successful"
// @Failure      400  {object}  map[string]string  "Invalid request"
// @Failure      401  {object}  map[string]string  "Invalid credentials"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /auth/login [post]
func LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	ctx := context.Background()
	params := auth_db.GetUserByIdentifyParams{
		Username: pgtype.Text{String: req.Identify, Valid: true},
		Email:    pgtype.Text{String: req.Identify, Valid: true},
		Phone:    pgtype.Text{String: req.Identify, Valid: true},
	}

	user, err := db.AuthQueries.GetUserByIdentify(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	t := jwt.New(jwt.SigningMethodHS256)
	jwtToken, err := t.SignedString([]byte("jwt"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"name": user.Name,
		},
		"token": jwtToken,
	})
}

// RegisterHandler godoc
// @Summary      User register
// @Description  Register a new user with username, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body RegisterRequest true "Register data"
// @Success      200 {object} map[string]string "{"success": "Register success"}"
// @Failure      400 {object} map[string]string "{"error": "invalid request body"}"
// @Failure      500 {object} map[string]string "{"error": "server error"}"
// @Router       /auth/register [post]
func RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest

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

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()
	if err := db.AuthQueries.CreateUser(ctx, auth_db.CreateUserParams{
		Name:     req.Name,
		Email:    pgtype.Text{String: req.Email, Valid: true},
		Phone:    pgtype.Text{String: req.Phone, Valid: true},
		Username: pgtype.Text{String: req.Username, Valid: true},
		Password: string(hash),
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"suceess": "Register success",
	})
}
