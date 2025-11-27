package customer

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// GetCustomersHandler godoc
// @Summary      Get customer list
// @Description  Returns a list of customers
// @Tags         customers
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  PaginatedResponse[any]
// @Router       /customers [get]
func GetCustomersHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()
	result, err := db.ProductQueries.GetCustomers(ctx, product_db.GetCustomersParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountCustomers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(PaginatedResponse[product_db.GetCustomersRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       result,
	})
}

// CreateCustomerHandler godoc
// @Summary      Create a new customer
// @Description  Creates a new customer and returns the created customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateCustomerRequest  true  "Create data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /customers [post]
func CreateCustomerHandler(c *fiber.Ctx) error {
	var req CreateCustomerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	customerParams := product_db.CreateCustomerParams{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: string(passwordHash),
	}

	customerID, err := db.ProductQueries.CreateCustomer(ctx, customerParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": customerID,
	})
}

// BulkDeleteCustomersHandler godoc
// @Summary      Delete multiple customer
// @Description  Deletes multiple customer by their IDs
// @Tags         customers
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteCustomersRequest  true  "List of order IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /customers [delete]
func BulkDeleteCustomersHandler(c *fiber.Ctx) error {
	var req DeleteCustomersRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx := context.Background()
	if err := db.ProductQueries.BulkDeleteCustomers(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

// RegisterCustomerHandler godoc
// @Summary      Create a new customer
// @Description  Creates a new customer and returns the created customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateCustomerRequest  true  "Create data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /customers/register [post]
func RegisterCustomerHandler(c *fiber.Ctx) error {
	var req CreateCustomerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	customerParams := product_db.CreateCustomerParams{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: string(passwordHash),
	}

	customerID, err := db.ProductQueries.CreateCustomer(ctx, customerParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": customerID,
	})
}

// CustomerLoginHandler godoc
// @Summary      Customer login
// @Description  Authenticates a customer using phone and password.
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        payload  body      CustomerLoginRequest  true  "Custoerm login request"
// @Success      200  {object}  map[string]interface{}  "Login successful"
// @Failure      400  {object}  map[string]string  "Invalid request"
// @Failure      401  {object}  map[string]string  "Invalid credentials"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /customers/login [post]
func CustomerLoginHandler(c *fiber.Ctx) error {
	var req CustomerLoginRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	ctx := context.Background()

	user, err := db.ProductQueries.GetCustomerByPhone(ctx, req.Phone)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"id":   user.ID,
			"name": user.Name,
		},
		"token": jwtToken,
	})
}
