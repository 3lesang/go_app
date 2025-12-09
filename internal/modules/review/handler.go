package review

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetReviewsByProductHandler godoc
// @Summary      Get review list
// @Description  Returns a list of reviews for a product, with optional filters
// @Tags         reviews
// @Produce      json
// @Param        id          path      int   true   "Product ID"
// @Param        page        query     int   false  "Page number"  default(1)
// @Param        page_size   query     int   false  "Page size"    default(10)
// @Param        rating      query     int   false  "Filter by rating"
// @Param        has_image   query     bool  false  "Filter by whether review has image (true/false)"
// @Param        sort_flag   query     int   false  "Sort by created_at (1 = newest, 0 = oldest)"
// @Success      200  {object}  map[string]interface{}
// @Router       /reviews/products/{id} [get]
func GetReviewsByProductHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	var rating pgtype.Int4
	if r := c.Query("rating"); r != "" {
		if val, err := strconv.Atoi(r); err == nil {
			rating = pgtype.Int4{Int32: int32(val), Valid: true}
		}
	}

	var hasImage pgtype.Bool
	if hi := c.Query("has_image"); hi != "" {
		if val, err := strconv.ParseBool(hi); err == nil {
			hasImage = pgtype.Bool{Bool: val, Valid: true}
		}
	}

	var sortFlag pgtype.Int4
	if s := c.Query("sort_flag"); s != "" {
		if val, err := strconv.Atoi(s); err == nil {
			sortFlag = pgtype.Int4{Int32: int32(val), Valid: true}
		}
	}

	ctx := context.Background()
	reviews, err := db.ProductQueries.GetReviewsByProductID(ctx, product_db.GetReviewsByProductIDParams{
		ProductID:   pgtype.Int8{Int64: id, Valid: true},
		LimitCount:  int32(pageSize),
		OffsetCount: int32(offset),
		HasFile:     hasImage,
		Rating:      rating,
		SortFlag:    sortFlag,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountReviewsByProduct(c.Context(), pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse[product_db.GetReviewsByProductIDRow]{
		Data:       reviews,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}

// CreateReviewHandler godoc
// @Summary      Create a new review
// @Description  Creates a new review with optional image and returns the created review
// @Tags         reviews
// @Security BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        product_id  formData  int64     true   "Product ID" default(0)
// @Param        customer_id formData  int64     true   "Customer ID" default(0)
// @Param        rating      formData  int     true   "Rating (1-5)" default(5)
// @Param        comment     formData  string  true   "Review comment" default(test)
// @Param 			 files 			 formData  []file   false   "Array of files to upload" collectionFormat multi
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /reviews [post]
func CreateReviewHandler(c *fiber.Ctx) error {
	var req CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form data"})
	}
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	files := form.File["files"]
	ctx := context.Background()

	params := product_db.CreateReviewParams{
		ProductID:  pgtype.Int8{Int64: req.ProductID, Valid: true},
		Rating:     pgtype.Int4{Int32: int32(req.Rating), Valid: true},
		Comment:    pgtype.Text{String: req.Comment, Valid: true},
		CustomerID: req.CustomerID,
		HasFile:    len(files) > 0,
	}

	id, err := db.ProductQueries.CreateReview(ctx, params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(files) > 0 {
		createReviewFileParams := product_db.BulkInsertReviewFilesParams{}
		for _, file := range files {
			fileKey := fmt.Sprintf("%s%s", "review/", file.Filename)
			createReviewFileParams.Names = append(createReviewFileParams.Names, fileKey)
			createReviewFileParams.ReviewIds = append(createReviewFileParams.ReviewIds, id)
		}
		db.ProductQueries.BulkInsertReviewFiles(ctx, createReviewFileParams)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      id,
		"message": "review created successfully",
	})
}

// GetOverviewByProductHandler godoc
// @Summary      Get average rating of a product
// @Description  Returns the average rating and total number of reviews for a specific product
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  AverageRatingResponse[any]
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /reviews/products/{id}/overview [get]
func GetOverviewByProductHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product id"})
	}

	ctx := c.Context()
	ratingInfo, err := db.ProductQueries.GetAverageRatingByProduct(ctx, pgtype.Int8{Int64: id, Valid: true})
	countFiles, err := db.ProductQueries.CountReviewFilesByProduct(ctx, pgtype.Int8{Int64: id, Valid: true})
	result, err := db.ProductQueries.GetReviewFilesByProduct(ctx, pgtype.Int8{Int64: id, Valid: true})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := AverageRatingResponse[product_db.GetReviewFilesByProductRow]{
		AverageRating: ratingInfo.AverageRating,
		TotalReviews:  ratingInfo.TotalReviews,
		TotalFiles:    countFiles,
		Data:          result,
	}
	return c.JSON(response)
}

// BulkDeleteReviewsHandler godoc
// @Summary      Delete multiple reviews
// @Description  Deletes multiple reviews by their IDs
// @Tags         reviews
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteReviewsRequest  true  "List of review IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /reviews [delete]
func BulkDeleteReviewsHandler(c *fiber.Ctx) error {
	var req DeleteReviewsRequest
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
	if err := db.ProductQueries.BulkDeleteReviews(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
