package order

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"encoding/json"
	"math"
	"strconv"

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
// @Success      200  {object}  PaginatedResponse[any]
// @Router       /orders [get]
func GetOrdersHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()
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

	items, err := db.ProductQueries.GetItemsByOrderID(ctx, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	orderItems := make([]OrderItemResponse, 0, len(items))
	for _, item := range items {
		var opts map[string]string
		if len(item.Options) > 0 {
			if err := json.Unmarshal(item.Options, &opts); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}

		orderItems = append(orderItems, OrderItemResponse{
			ID:        item.OrderItemID,
			Quantity:  item.Quantity,
			SalePrice: item.SalePrice,
			ProductID: item.ProductID,
			Name:      item.Name,
			Slug:      item.Slug,
			Options:   opts,
		})
	}

	return c.Status(fiber.StatusOK).JSON(OrderResponse{
		FullName:       result.FullName.String,
		Phone:          result.Phone.String,
		AddressLine:    result.AddressLine.String,
		TotalAmount:    result.TotalAmount,
		DiscountAmount: result.DiscountAmount,
		Items:          orderItems,
	})
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
		TotalAmount:    req.TotalAmount,
		DiscountAmount: req.DiscountAmount,
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
		createOrderItemParams.VariantIds = append(createOrderItemParams.VariantIds, item.VariantID)
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
