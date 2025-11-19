package discount

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateDiscountHandler godoc
// @Summary Create a new discount
// @Description Creates a discount
// @Tags discounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CreateDiscountRequest true "Discount data"
// @Success 201 {object} int64
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts [post]
func CreateDiscountHandler(c *fiber.Ctx) error {
	var req CreateDiscountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	discountID, err := db.ProductQueries.CreateDiscount(ctx,
		product_db.CreateDiscountParams{
			Code:             pgtype.Text{String: req.Code, Valid: len(req.Code) > 0},
			Title:            req.Title,
			DiscountType:     req.DiscountType,
			Status:           req.Status,
			UsageLimit:       pgtype.Int4{Int32: req.UsageLimit, Valid: req.UsageLimit > 0},
			PerCustomerLimit: pgtype.Int4{Int32: req.PerCustomerLimit, Valid: req.PerCustomerLimit > 0},
			StartsAt:         pgtype.Timestamp{Time: req.StartsAt, Valid: true},
			EndsAt:           pgtype.Timestamp{Time: req.EndsAt, Valid: true},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(discountID)
}

// CreateDiscountTargetHandler godoc
// @Summary Create discount target
// @Description Creates discount target
// @Tags discounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CreateDiscountTargetRequest true "Discount data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts/{id}/targets [post]
func CreateDiscountTargetHandler(c *fiber.Ctx) error {
	var req CreateDiscountTargetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	idParam := c.Params("id")
	discountID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid param"})
	}

	createDiscountTarget := product_db.BulkInsertDiscountTargetsParams{}

	for _, id := range req.IDs {
		createDiscountTarget.TargetTypes = append(createDiscountTarget.TargetTypes, req.TargetType)
		createDiscountTarget.DiscountIds = append(createDiscountTarget.DiscountIds, discountID)
		createDiscountTarget.TargetIds = append(createDiscountTarget.TargetIds, id)
	}
	return c.SendStatus(fiber.StatusCreated)
}

// CreateDiscountEffectHandler godoc
// @Summary Create discount effect
// @Description Creates discount effect
// @Tags discounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CreateDiscountEffectRequest true "Discount data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts/{id}/effects [post]
func CreateDiscountEffectHandler(c *fiber.Ctx) error {
	var req CreateDiscountEffectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	idParam := c.Params("id")
	discountID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid param"})
	}

	createDiscountTarget := product_db.CreateDiscountEffectParams{
		DiscountID: discountID,
		EffectType: req.EffectType,
		Value:      pgtype.Text{String: req.Value, Valid: true},
		AppliesTo:  req.AppliesTo,
	}
	ctx := context.Background()
	db.ProductQueries.CreateDiscountEffect(ctx, createDiscountTarget)
	return c.SendStatus(fiber.StatusCreated)
}

// CreateDiscountConditionHandler godoc
// @Summary Create discount condition
// @Description Creates discount condition
// @Tags discounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CreateDiscountConditionRequest true "Discount data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts/{id}/conditions [post]
func CreateDiscountConditionHandler(c *fiber.Ctx) error {
	var req CreateDiscountConditionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	idParam := c.Params("id")
	discountID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid param"})
	}

	params := product_db.CreateDiscountConditionParams{
		DiscountID:    discountID,
		ConditionType: req.ConditionType,
		Operator:      req.Operator,
		Value:         req.Value,
	}
	ctx := context.Background()
	db.ProductQueries.CreateDiscountCondition(ctx, params)
	return c.SendStatus(fiber.StatusCreated)
}

// UpdateDiscountHandler godoc
// @Summary Update a discount
// @Description Updates a discount by ID
// @Tags discounts
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path int true "Discount ID"
// @Param payload body UpdateDiscountRequest true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts/{id} [put]
func UpdateDiscountHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid param"})
	}

	var req UpdateDiscountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	discount, err := db.ProductQueries.UpdateDiscount(ctx, product_db.UpdateDiscountParams{
		ID:               id,
		Title:            req.Title,
		Status:           req.Status,
		UsageLimit:       pgtype.Int4{Int32: req.UsageLimit, Valid: req.UsageLimit > 0},
		PerCustomerLimit: pgtype.Int4{Int32: req.PerCustomerLimit, Valid: req.PerCustomerLimit > 0},
		StartsAt:         pgtype.Timestamp{Time: req.StartsAt, Valid: true},
		EndsAt:           pgtype.Timestamp{Time: req.EndsAt, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(discount)
}

// UpdateDiscountEffectHandler godoc
// @Summary Update a discount effect
// @Description Updates a discount effect by ID
// @Tags discounts
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path int true "Discount ID"
// @Param payload body UpdateDiscountRequest true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts/{id}/effects/{effectID} [put]
func UpdateDiscountEffectHandler(c *fiber.Ctx) error {
	idParam := c.Params("effectID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid param"})
	}

	var req UpdateDiscountEffectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	err = db.ProductQueries.UpdateDiscountEffect(ctx, product_db.UpdateDiscountEffectParams{
		ID:         id,
		EffectType: req.EffectType,
		Value:      pgtype.Text{String: req.Value, Valid: len(req.Value) > 0},
		AppliesTo:  req.AppliesTo,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdateDiscountConditionHandler godoc
// @Summary Update a discount condition
// @Description Updates a discount condition by ID
// @Tags discounts
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path int true "Discount ID"
// @Param payload body UpdateDiscountRequest true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts/{id}/conditions/{conditionID} [put]
func UpdateDiscountConditionHandler(c *fiber.Ctx) error {
	idParam := c.Params("conditionID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid param"})
	}

	var req UpdateDiscountConditionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	err = db.ProductQueries.UpdateDiscountCondition(ctx, product_db.UpdateDiscountConditionParams{
		ID:            id,
		ConditionType: req.ConditionType,
		Value:         req.Value,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetDiscountsHandler godoc
// @Summary      Get discounts list
// @Description  Returns a paginated list of discounts
// @Tags         discounts
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        page_size query     int  false  "Page size"    default(10)
// @Success      200  {object}  PaginatedResponse[any]
// @Failure      500  {object}  map[string]string
// @Router       /discounts [get]
func GetDiscountsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	ctx := context.Background()

	discounts, err := db.ProductQueries.ListDiscounts(ctx, product_db.ListDiscountsParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	total, err := db.ProductQueries.CountDiscounts(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(PaginatedResponse[product_db.Discount]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(total),
		TotalPages: totalPages,
		Data:       discounts,
	})
}

// GetDiscountHandler godoc
// @Summary Get a discount by ID
// @Description Returns a discount by ID
// @Tags discounts
// @Produce json
// @Param id path int true "Discount ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts/{id} [get]
func GetDiscountHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid param"})
	}

	ctx := context.Background()
	discount, err := db.ProductQueries.GetDiscountWithRelations(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if discount.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Discount not found"})
	}

	return c.JSON(discount)
}

// BulkDeleteDiscountsHandler godoc
// @Summary Bulk delete discounts
// @Description Deletes multiple discounts by IDs
// @Tags discounts
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param payload body BulkDeleteDiscountsRequest true "IDs to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /discounts [delete]
func BulkDeleteDiscountsHandler(c *fiber.Ctx) error {
	var req BulkDeleteDiscountsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	if err := db.ProductQueries.BulkDeleteDiscounts(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted": len(req.IDs)})
}
