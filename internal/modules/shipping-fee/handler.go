package shippingfee

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetShippingFeesHandler godoc
// @Summary      Get shipping fee list
// @Description  Returns a list of shipping fee
// @Tags         shipping-fees
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  map[string]interface{}
// @Router       /shipping-fees [get]
func GetShippingFeesHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()
	result, err := db.ProductQueries.GetShippingFees(ctx, product_db.GetShippingFeesParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountShippingFees(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse[product_db.GetShippingFeesRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       result,
	})
}

// GetShippingFeeHandler godoc
// @Summary      Get a shipping fee
// @Description  Returns a shipping fee by ID
// @Tags         shipping-fees
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /shipping-fees/{id} [get]
func GetShippingFeeHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	result, err := db.ProductQueries.GetShippingFee(ctx, id)
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

	return c.Status(fiber.StatusOK).JSON(ShippingFeeResponse{
		ID:            result.ID,
		MinWeight:     result.MinWeight,
		MaxWeight:     result.MaxWeight,
		FeeAmount:     result.FeeAmount,
		MinOrderValue: result.MinOrderValue.Int32,
		FreeShipping:  result.FreeShipping.Bool,
	})
}

// GetShippingFeeByWeightHandler godoc
// @Summary      Get a shipping fee
// @Description  Returns a shipping fee by weight
// @Tags         shipping-fees
// @Security BearerAuth
// @Produce      json
// @Param        value   path      int  true  "value"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /shipping-fees/weight/{value} [get]
func GetShippingFeeByWeightHandler(c *fiber.Ctx) error {
	param := c.Params("value")
	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	result, err := db.ProductQueries.GetShippingFeeByWeight(ctx, int32(value))
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

	return c.Status(fiber.StatusOK).JSON(result)
}

// CreateShippingFeeHandler godoc
// @Summary      Create a new shipping fee
// @Description  Creates a new shipping fee and returns id
// @Tags         shipping-fees
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateShippingFeeRequest  true  "Create shpping fee data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /shipping-fees [post]
func CreateShippingFeeHandler(c *fiber.Ctx) error {
	var req CreateShippingFeeRequest

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
	params := product_db.CreateShippingFeeParams{
		MinWeight:     req.MinWeight,
		MaxWeight:     req.MaxWeight,
		FeeAmount:     req.FeeAmount,
		MinOrderValue: pgtype.Int4{Int32: req.MinOrderValue, Valid: true},
		FreeShipping:  pgtype.Bool{Bool: req.FreeShipping, Valid: true},
	}
	idReturn, err := db.ProductQueries.CreateShippingFee(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(idReturn)
}

// UpdateShippingFeeHandler godoc
// @Summary      Update a shipping fee
// @Description  Updates a shipping fee
// @Tags         shipping-fees
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Param        payload  body	UpdateShippingFeeRequest  true  "Update shipping fee data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /shipping-fees/{id} [put]
func UpdateShippingFeeHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	var req UpdateShippingFeeRequest
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
	params := product_db.UpdateShippingFeeParams{
		ID:            id,
		MinWeight:     req.MinWeight,
		MaxWeight:     req.MaxWeight,
		FeeAmount:     req.FeeAmount,
		MinOrderValue: pgtype.Int4{Int32: req.MinOrderValue, Valid: true},
		FreeShipping:  pgtype.Bool{Bool: req.FreeShipping, Valid: true},
	}
	if err := db.ProductQueries.UpdateShippingFee(ctx, params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

// DeleteShippingFeesHandler godoc
// @Summary      Delete multiple shipping fees
// @Description  Deletes multiple shipping fees by their IDs
// @Tags         shipping-fees
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteShippingFeesRequest  true  "List of shipping fee IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /shipping-fees [delete]
func DeleteShippingFeesHandler(c *fiber.Ctx) error {
	var req DeleteShippingFeesRequest
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
	if err := db.ProductQueries.BulkDeleteShippingFees(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
