package file

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetFilesHandler godoc
// @Summary      Get file list
// @Description  Returns a list of files
// @Tags         files
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  PaginatedResponse[any]
// @Router       /files [get]
func GetFilesHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()

	files, err := db.ProductQueries.GetFiles(ctx, product_db.GetFilesParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountFiles(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(PaginatedResponse[product_db.GetFilesRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       files,
	})
}

// CreateFileHandler godoc
// @Summary      Create a new file
// @Description  Creates a new file and returns the created file
// @Tags         files
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateFileRequest  true  "Create data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /files [post]
func CreateFileHandler(c *fiber.Ctx) error {
	var req CreateFileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	params := []string{}
	for _, file := range req.Names {
		params = append(params, file)
	}
	ctx := context.Background()
	if err := db.ProductQueries.BulkInsertFiles(ctx, params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(params)
}

// DeleteFilesHandler godoc
// @Summary      Delete multiple files
// @Description  Deletes multiple files by their IDs
// @Tags         files
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteFilesRequest  true  "List of file IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /files [delete]
func DeleteFilesHandler(c *fiber.Ctx) error {
	var req DeleteFilesRequest
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
	if err := db.ProductQueries.BulkDeleteFiles(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "users deleted successfully",
		"count":   len(req.IDs),
	})
}
