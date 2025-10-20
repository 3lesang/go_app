package product

import (
	database "app/internal/database/postgres"
	"context"
	"database/sql"
	"math/big"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetProductsHandler godoc
// @Summary      Get product list
// @Description  Returns a list of products
// @Tags         products
// @Security BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /products [get]
func GetProductsHandler(c *fiber.Ctx) error {
	ctx := context.Background()
	products, err := database.PGQueries.ListProducts(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": products,
	})
}

// GetProductHandler godoc
// @Summary      Get a product
// @Description  Returns a product by ID
// @Tags         products
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [get]
func GetProductHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	product, err := database.PGQueries.GetProduct(ctx, id)
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
	return c.Status(fiber.StatusOK).JSON(product)
}

// CreateProductHandler godoc
// @Summary      Create a new product
// @Description  Creates a new product and returns the created product
// @Tags         products
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        product  body	CreateProductRequest  true  "Product data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products [post]
func CreateProductHandler(c *fiber.Ctx) error {
	var req CreateProductRequest

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
	if err := database.PGQueries.CreateProduct(ctx, database.CreateProductParams{
		Name:        req.Name,
		Slug:        pgtype.Text{String: req.Slug, Valid: true},
		OriginPrice: pgtype.Numeric{Int: big.NewInt(req.OriginPrice), Valid: true},
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": "Create product successfully",
	})
}

// UpdateProductHandler godoc
// @Summary      Update a product
// @Description  Updates a product
// @Tags         products
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Param        product  body	UpdateProductRequest  true  "Product data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [put]
func UpdateProductHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	var req UpdateProductRequest
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
	if err := database.PGQueries.UpdateProduct(ctx, database.UpdateProductParams{
		Name: req.Name,
		Slug: pgtype.Text{String: req.Slug, Valid: true},
		ID:   id,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteProductsHandler godoc
// @Summary      Delete multiple products
// @Description  Deletes multiple products by their IDs
// @Tags         products
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteProductsRequest  true  "List of product IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products [delete]
func DeleteProductsHandler(c *fiber.Ctx) error {
	var req DeleteProductsRequest
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
	if err := database.PGQueries.DeleteProducts(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
