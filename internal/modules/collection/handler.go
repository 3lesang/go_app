package collection

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"database/sql"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetCollectionsHandler godoc
// @Summary      Get collection list
// @Description  Returns a list of collections
// @Tags         collections
// @Security BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /collections [get]
func GetCollectionsHandler(c *fiber.Ctx) error {
	ctx := context.Background()
	results, err := db.ProductQueries.GetCollections(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": results,
	})
}

// GetCollectionHandler godoc
// @Summary      Get a collection
// @Description  Returns a collection by ID
// @Tags         collections
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /collections/{id} [get]
func GetCollectionHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	result, err := db.ProductQueries.GetCollection(ctx, id)
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

	products, err := db.ProductQueries.GetProductsByCollectionID(ctx, result.ID)

	return c.Status(fiber.StatusOK).JSON(CollectionResponse{
		ID:              result.ID,
		File:            result.File.String,
		Name:            result.Name,
		Slug:            result.Slug,
		MetaTitle:       result.MetaTitle.String,
		MetaDescription: result.MetaDescription.String,
		Products:        products,
	})
}

// CreateCollectionHandler godoc
// @Summary      Create a new collection
// @Description  Creates a new collection and returns the created collection
// @Tags         collections
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateCollectionRequest  true  "Create collection data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /collections [post]
func CreateCollectionHandler(c *fiber.Ctx) error {
	var req CreateCollectionRequest

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
	params := product_db.CreateCollectionParams{
		Name: req.Name,
		Slug: req.Slug,
		Layout: pgtype.Text{
			String: req.Layout,
			Valid:  true,
		},
		File:      pgtype.Text{String: req.File, Valid: req.File != ""},
		MetaTitle: pgtype.Text{String: req.MetaTitle, Valid: true},
		MetaDescription: pgtype.Text{
			String: req.MetaDescription,
			Valid:  true,
		},
	}
	collectionID, err := db.ProductQueries.CreateCollection(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	createProductCollectionParams := product_db.BulkInsertProductCollectionParams{}

	for _, id := range req.ProductIDs {
		createProductCollectionParams.CollectionIds = append(createProductCollectionParams.CollectionIds, collectionID)
		createProductCollectionParams.ProductIds = append(createProductCollectionParams.ProductIds, id)
	}
	db.ProductQueries.BulkInsertProductCollection(ctx, createProductCollectionParams)
	return c.SendStatus(fiber.StatusCreated)
}

// UpdateCollectionHandler godoc
// @Summary      Update a collection
// @Description  Updates a collection
// @Tags         collections
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Param        payload  body	UpdateCollectionRequest  true  "Update data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /collections/{id} [put]
func UpdateCollectionHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	var req UpdateCollectionRequest
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
	params := product_db.UpdateCollectionParams{
		ID:   id,
		Name: req.Name,
		Slug: req.Slug,
		File: pgtype.Text{String: req.File, Valid: true},
	}
	if err := db.ProductQueries.UpdateCollection(ctx, params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	insertCollectionProductParams := product_db.BulkInsertProductCollectionParams{}

	for _, pid := range req.ProductIDs {
		insertCollectionProductParams.CollectionIds = append(insertCollectionProductParams.CollectionIds, id)
		insertCollectionProductParams.ProductIds = append(insertCollectionProductParams.ProductIds, pid)
	}

	db.ProductQueries.DeleteProductsByCollectionID(ctx, id)
	db.ProductQueries.BulkInsertProductCollection(ctx, insertCollectionProductParams)
	return c.SendStatus(fiber.StatusOK)
}

// DeleteCollectionsHandler godoc
// @Summary      Delete multiple collections
// @Description  Deletes multiple collections by their IDs
// @Tags         collections
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteCollectionsRequest  true  "List of IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /collections [delete]
func DeleteCollectionsHandler(c *fiber.Ctx) error {
	var req DeleteCollectionsRequest
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
	if err := db.ProductQueries.BulkDeleteCollections(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
