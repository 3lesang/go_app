package order

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetOrdersHandler godoc
// @Summary      Get order list
// @Description  Returns a list of orders
// @Tags         orders
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Param        status query     string     false  "Status" Enums(pending, confirmed, shipping, shipped, canceled)    default(”)
// @Success      200  {object}  PaginatedResponse[any]
// @Router       /orders [get]
func GetOrdersHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status", "")

	offset := (page - 1) * pageSize

	ctx := context.Background()
	if len(status) > 0 {
		result, err := db.ProductQueries.GetOrdersByStatus(ctx, product_db.GetOrdersByStatusParams{
			Limit:  int32(pageSize),
			Offset: int32(offset),
			Status: pgtype.Text{String: status, Valid: true},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		total, err := db.ProductQueries.CountOrdersByStatus(ctx, pgtype.Text{String: status, Valid: true})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		return c.JSON(PaginatedResponse[product_db.GetOrdersByStatusRow]{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
			Data:       result,
		})
	}

	result, err := db.ProductQueries.GetOrders(ctx, product_db.GetOrdersParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountOrders(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(PaginatedResponse[product_db.GetOrdersRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       result,
	})
}

// CountOrderHandler godoc
// @Summary      Get count status order
// @Description  Returns count of status order
// @Tags         orders
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /count-orders [get]
func CountOrderHandler(c *fiber.Ctx) error {
	ctx := context.Background()
	result, err := db.ProductQueries.CountStatusOrder(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// GetOrderHandler godoc
// @Summary      Get a order
// @Description  Returns a order by ID
// @Tags         orders
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders/{id} [get]
func GetOrderHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	result, err := db.ProductQueries.GetOrder(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// CheckOrderCreatedHandler godoc
// @Summary      Check order is created success
// @Description  Returns a order by ID
// @Tags         orders
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders/{id}/success [get]
func CheckOrderCreatedHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	id, err = db.ProductQueries.CheckOrderCreated(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func generateCode() string {
	datePart := time.Now().Format("20060102")

	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())

	suffix := make([]byte, 5)
	for i := range suffix {
		suffix[i] = letters[rand.Intn(len(letters))]
	}

	return datePart + string(suffix)
}

// CreateOrderHandler godoc
// @Summary      Create a new order
// @Description  Creates a new order and returns the created order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateOrderRequest  true  "Create data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders [post]
func CreateOrderHandler(c *fiber.Ctx) error {
	var req CreateOrderRequest

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
	addressParams := product_db.CreateAddressParams{
		FullName:    req.Address.FullName,
		Email:       pgtype.Text{String: req.Address.Email, Valid: true},
		Phone:       pgtype.Text{String: req.Address.Phone, Valid: true},
		AddressLine: req.Address.AddressLine,
	}
	addressID, err := db.ProductQueries.CreateAddress(ctx, addressParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	orderParams := product_db.CreateOrderParams{
		Code:           generateCode(),
		TotalAmount:    req.TotalAmount,
		DiscountAmount: req.DiscountAmount,
		ShippingFeeAmount: pgtype.Int4{
			Int32: req.ShippingFeeAmount,
			Valid: true,
		},
		ShippingAddressID: pgtype.Int8{
			Int64: addressID,
			Valid: true,
		},
	}
	orderID, err := db.ProductQueries.CreateOrder(ctx, orderParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	createOrderItemParams := product_db.BulkInsertOrderItemsParams{}
	for _, item := range req.Items {
		createOrderItemParams.Quantities = append(createOrderItemParams.Quantities, item.Quantity)
		createOrderItemParams.SalePrices = append(createOrderItemParams.SalePrices, item.SalePrice)
		createOrderItemParams.OrderIds = append(createOrderItemParams.OrderIds, orderID)
		createOrderItemParams.ProductIds = append(createOrderItemParams.ProductIds, item.ProductID)
		if item.VariantID != 0 {
			createOrderItemParams.VariantIds = append(createOrderItemParams.VariantIds, item.VariantID)
		}
	}
	if err := db.ProductQueries.BulkInsertOrderItems(ctx, createOrderItemParams); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": orderID,
	})
}

// UpdateOrderStatusHandler godoc
// @Summary      Update status order
// @Description  Update status order by id
// @Tags         orders
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders/{id}/status [put]
func UpdateOrderStatusHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	var req UpdateOrderRequest
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
	if err := db.ProductQueries.UpdateOrder(ctx, product_db.UpdateOrderParams{
		ID: id,
		Status: pgtype.Text{
			String: req.Status,
			Valid:  true,
		},
		CancelReason: pgtype.Text{String: req.CancelReason, Valid: true},
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

// DeleteOrdersHandler godoc
// @Summary      Delete multiple orders
// @Description  Deletes multiple orders by their IDs
// @Tags         orders
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteOrdersRequest  true  "List of order IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders [delete]
func DeleteOrdersHandler(c *fiber.Ctx) error {
	var req DeleteOrdersRequest
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
	if err := db.ProductQueries.BulkDeleteOrders(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
