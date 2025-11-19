package category

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// GetCategoriesHandler godoc
// @Summary      Get category list
// @Description  Returns a list of categories
// @Tags         categories
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  map[string]interface{}
// @Router       /categories [get]
func GetCategoriesHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()
	result, err := db.ProductQueries.GetCategories(ctx, product_db.GetCategoriesParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountCategories(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse[product_db.GetCategoriesRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       result,
	})
}

// GetCategoryHandler godoc
// @Summary      Get a category
// @Description  Returns a category by ID
// @Tags         categories
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories/{id} [get]
func GetCategoryHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	product, err := db.ProductQueries.GetCategory(ctx, id)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":   product.ID,
		"name": product.Name,
		"slug": product.Slug,
	})
}

// CreateCategoryHandler godoc
// @Summary      Create a new category
// @Description  Creates a new category and returns the created category
// @Tags         categories
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateCategoryRequest  true  "Create data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories [post]
func CreateCategoryHandler(c *fiber.Ctx) error {
	var req CreateCategoryRequest

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
	params := product_db.CreateCategoryParams{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := db.ProductQueries.CreateCategory(ctx, params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": "Create successfully",
	})
}

// UpdateCategoryHandler godoc
// @Summary      Update a category
// @Description  Updates a category
// @Tags         categories
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Param        payload  body	UpdateCategoryRequest  true  "Update data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories/{id} [put]
func UpdateCategoryHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	var req UpdateCategoryRequest
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
	params := product_db.UpdateCategoryParams{
		ID:   id,
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := db.ProductQueries.UpdateCategory(ctx, params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteCategoriesHandler godoc
// @Summary      Delete multiple categories
// @Description  Deletes multiple categories by their IDs
// @Tags         categories
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteCategoriesRequest  true  "List of IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories [delete]
func DeleteCategoriesHandler(c *fiber.Ctx) error {
	var req DeleteCategoriesRequest
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
	if err := db.ProductQueries.BulkDeleteCategories(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
