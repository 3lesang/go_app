package post

import (
	"app/internal/db"
	blog_db "app/internal/db/blog"
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// GetPostsHandler godoc
// @Summary      Get post list
// @Description  Returns a list of posts
// @Tags         posts
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  PaginatedResponse[any]
// @Router       /posts [get]
func GetPostsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()
	result, err := db.BlogQueries.GetPosts(ctx, blog_db.GetPostsParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.BlogQueries.CountPosts(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(PaginatedResponse[blog_db.GetPostsRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       result,
	})
}

// GetPostHandler godoc
// @Summary      Get a post
// @Description  Returns a post by ID
// @Tags         posts
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{id} [get]
func GetPostHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	result, err := db.BlogQueries.GetPost(ctx, id)
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

// GetPostBySlugHandler godoc
// @Summary      Get a post
// @Description  Returns a post by ID
// @Tags         posts
// @Security BearerAuth
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/slug/{id} [get]
func GetPostBySlugHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	ctx := context.Background()
	result, err := db.BlogQueries.GetPostBySlug(ctx, param)
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

// CreatePostHandler godoc
// @Summary      Create a new post
// @Description  Creates a new post and returns the created post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        payload  body	CreatePostRequest  true  "Create data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts [post]
func CreatePostHandler(c *fiber.Ctx) error {
	var req CreatePostRequest

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

	postParams := blog_db.CreatePostParams{
		Title: req.Title,
		Slug:  req.Slug,
	}

	postID, err := db.BlogQueries.CreatePost(ctx, postParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": postID,
	})
}

// UpdatePostHandler godoc
// @Summary      Update a post
// @Description  Updates a post
// @Tags         posts
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Param        payload  body	UpdatePostRequest  true  "Update data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts/{id} [put]
func UpdatePostHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	var req UpdatePostRequest
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
	params := blog_db.UpdatePostParams{
		ID:    id,
		Title: req.Title,
		Slug:  req.Slug,
	}
	if err := db.BlogQueries.UpdatePost(ctx, params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// BulkDeletePostsHandler godoc
// @Summary      Delete multiple post
// @Description  Deletes multiple post by their IDs
// @Tags         posts
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeletePostsRequest  true  "List of order IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /posts [delete]
func BulkDeletePostsHandler(c *fiber.Ctx) error {
	var req DeletePostsRequest
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
	if err := db.BlogQueries.BulkDeletePosts(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
