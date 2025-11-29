package hotspot

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetHotspotsHandler godoc
// @Summary      Get hotspot list
// @Description  Returns a list of hotspots
// @Tags         hotspots
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  PaginatedResponse[any]
// @Router       /hotspots [get]
func GetHotspotsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()

	files, err := db.ProductQueries.GetHotspots(ctx, product_db.GetHotspotsParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountHotspots(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(PaginatedResponse[product_db.Hotspot]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       files,
	})
}

// GetHotspotHandler godoc
// @Summary      Get a hotspot
// @Description  Returns a hotspot by ID
// @Tags         hotspots
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hotspots/{id} [get]
func GetHotspotHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	result, err := db.ProductQueries.GetHotspot(ctx, id)
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

// CreateHotspotHandler godoc
// @Summary      Create a new hotspot
// @Description  Creates a new hotspot and returns the created hotspot
// @Tags         hotspots
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateHotspotRequest  true  "Create data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hotspots [post]
func CreateHotspotHandler(c *fiber.Ctx) error {
	var req CreateHotspotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	ctx := context.Background()
	hotspotID, err := db.ProductQueries.CreateHotspot(ctx, req.File)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(req.Spots) > 0 {
		params := product_db.BulkInsertProductHotspotsParams{}
		for _, spot := range req.Spots {
			params.ProductIds = append(params.ProductIds, spot.ProductID)
			params.HotspotIds = append(params.HotspotIds, hotspotID)
			params.Xs = append(params.Xs, spot.X)
			params.Ys = append(params.Ys, spot.Y)
		}
		if err := db.ProductQueries.BulkInsertProductHotspots(ctx, params); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(hotspotID)
}

func ToSet[T comparable](slice []T) map[T]struct{} {
	set := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		set[v] = struct{}{}
	}
	return set
}

// UpdateHotspotHandler godoc
// @Summary      Update a hotspot
// @Description  Updates a hotspot
// @Tags         hotspots
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Param        payload  body	UpdateHotspotRequest  true  "Update data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hotspots/{id} [put]
func UpdateHotspotHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}

	var req UpdateHotspotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	productHotspotIDs, err := db.ProductQueries.GetProductHotspotsByHotspot(ctx, id)
	if err != nil {
		return err
	}

	existingSet := ToSet(productHotspotIDs)

	createProductHotspot := product_db.BulkInsertProductHotspotsParams{}
	updateProductHotspot := product_db.BulkUpdateProductHotspotsParams{}

	incomingSet := make(map[int64]struct{}, len(req.Spots))

	for _, spot := range req.Spots {
		incomingSet[spot.ID] = struct{}{}

		if _, exists := existingSet[spot.ID]; exists {
			updateProductHotspot.ProductIds = append(updateProductHotspot.ProductIds, spot.ProductID)
			updateProductHotspot.HotspotIds = append(updateProductHotspot.HotspotIds, id)
		} else {
			createProductHotspot.HotspotIds = append(createProductHotspot.HotspotIds, id)
			createProductHotspot.ProductIds = append(createProductHotspot.ProductIds, spot.ProductID)
			createProductHotspot.Xs = append(createProductHotspot.Xs, spot.X)
			createProductHotspot.Ys = append(createProductHotspot.Ys, spot.Y)
		}
	}

	var deleteIDs []int64
	for _, existingID := range productHotspotIDs {
		if _, exists := incomingSet[existingID]; !exists {
			deleteIDs = append(deleteIDs, existingID)
		}
	}

	if len(createProductHotspot.HotspotIds) > 0 {
		err := db.ProductQueries.BulkInsertProductHotspots(ctx, createProductHotspot)
		if err != nil {
			return err
		}
	}

	if len(updateProductHotspot.HotspotIds) > 0 {
		err := db.ProductQueries.BulkUpdateProductHotspots(ctx, updateProductHotspot)
		if err != nil {
			return err
		}
	}

	if len(deleteIDs) > 0 {
		err := db.ProductQueries.BulkDeleteProductHotspots(ctx, deleteIDs)
		if err != nil {
			return err
		}
	}

	return c.SendStatus(fiber.StatusOK)

}

// DeleteHotspotsHandler godoc
// @Summary      Delete multiple hotspots
// @Description  Deletes multiple hotspots by their IDs
// @Tags         hotspots
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteHotspotsRequest  true  "List of file IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hotspots [delete]
func DeleteHotspotsHandler(c *fiber.Ctx) error {
	var req DeleteHotspotsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no user IDs provided",
		})
	}

	ctx := context.Background()
	if err := db.ProductQueries.BulkDeleteHotspots(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "users deleted successfully",
		"count":   len(req.IDs),
	})
}
