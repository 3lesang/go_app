package auth

import (
	"app/internal/database"
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// SignInHandler godoc
// @Summary      User Sign In
// @Description  Authenticates a user using email or username and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      SignInParams  true  "Sign In parameters"
// @Success      200  {object}  map[string]interface{}  "Sign in successful"
// @Failure      400  {object}  map[string]string  "Invalid request"
// @Failure      401  {object}  map[string]string  "Invalid credentials"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /auth/signin [post]
func SignInHandler(c *fiber.Ctx) error {
	var req SignInParams
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	ctx := context.Background()
	params := database.GetUserByIdentifyParams{
		Username: &req.Identify,
		Email:    &req.Identify,
		Phone:    &req.Identify,
	}

	user, err := database.SQliteQueries.GetUserByIdentify(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server error",
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

// SignUpHandler godoc
// @Summary      User Sign Up
// @Description  Register a new user with username, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        signup body SignUpParams true "Sign Up Params"
// @Success      200 {object} map[string]string "{"success": "sign up success"}"
// @Failure      400 {object} map[string]string "{"error": "invalid request body"}"
// @Failure      500 {object} map[string]string "{"error": "server error"}"
// @Router       /auth/signup [post]
func SignUpHandler(c *fiber.Ctx) error {
	var req SignUpParams

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
	if err := database.SQliteQueries.CreateUser(ctx, database.CreateUserParams{
		Name:     req.Name,
		Email:    &req.Email,
		Phone:    &req.Phone,
		Username: &req.Username,
		Password: string(hash),
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"suceess": "sign up success",
	})
}
