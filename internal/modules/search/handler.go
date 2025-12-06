package search

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"math"

	"github.com/gofiber/fiber/v2"
)

// SearchProductsHandler godoc
// @Summary      Search products
// @Description  Search products by keyword with pagination
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        keyword  query     string  false  "Search keyword"
// @Param        page     query     int     false  "Page number"      default(1)
// @Param        limit    query     int     false  "Items per page"   default(20)
// @Success      200      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /search [get]
func SearchProductsHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	keyword := c.Query("keyword", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	params := product_db.SearchProductsParams{
		Column1: keyword,
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	products, err := db.ProductQueries.SearchProducts(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	total, err := db.ProductQueries.CountSearchProducts(ctx, keyword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(PaginatedResponse[product_db.SearchProductsRow]{
		Data:       products,
		TotalItems: total,
		Page:       page,
		PageSize:   limit,
		TotalPages: totalPages,
	})
}
