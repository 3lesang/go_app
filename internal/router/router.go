package router

import (
	"app/internal/modules/auth"
	"app/internal/modules/category"
	"app/internal/modules/collection"
	"app/internal/modules/customer"
	"app/internal/modules/file"
	"app/internal/modules/order"
	"app/internal/modules/product"
	"app/internal/modules/review"
	"app/internal/modules/user"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	v1 := app.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authGroup.Post("/login", auth.LoginHandler)
	authGroup.Post("/register", auth.RegisterHandler)

	userGroup := v1.Group("/users")
	userGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	userGroup.Get("/", user.GetUsersHandler)
	userGroup.Get("/:id", user.GetUserHandler)
	userGroup.Post("/", user.CreateUserHandler)
	userGroup.Delete("/", user.DeleteUsersHandler)

	productGroup := v1.Group("/products")
	productGroup.Get("/slug/:slug", product.GetProductBySlugHandler)

	productGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	productGroup.Get("/", product.GetProductsHandler)
	productGroup.Get("/:id", product.GetProductHandler)
	productGroup.Post("/", product.CreateProductHandler)
	productGroup.Put("/:id", product.UpdateProductHandler)
	productGroup.Delete("/", product.DeleteProductsHandler)

	categoryGroup := v1.Group("/categories")
	categoryGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	categoryGroup.Get("/", category.GetCategoriesHandler)
	categoryGroup.Get("/:id", category.GetCategoryHandler)
	categoryGroup.Post("/", category.CreateCategoryHandler)
	categoryGroup.Put("/:id", category.UpdateCategoryHandler)
	categoryGroup.Delete("/", category.DeleteCategoriesHandler)

	collectionGroup := v1.Group("/collections")
	collectionGroup.Get("/hero", collection.GetHeroCollectionsHandler)
	collectionGroup.Get("/home", collection.GetHomeCollectionsHandler)
	collectionGroup.Get("/slug/:slug", collection.GetCollectionBySlugHandler)
	collectionGroup.Get("/:id/products", collection.GetProductsHandler)

	collectionGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	collectionGroup.Get("/", collection.GetCollectionsHandler)
	collectionGroup.Get("/:id", collection.GetCollectionHandler)
	collectionGroup.Post("/", collection.CreateCollectionHandler)
	collectionGroup.Put("/:id", collection.UpdateCollectionHandler)
	collectionGroup.Delete("/", collection.DeleteCollectionsHandler)

	orderGroup := v1.Group("/orders")
	orderGroup.Get("/", order.GetOrdersHandler)
	orderGroup.Get("/:id", order.GetOrderHandler)
	orderGroup.Get("/:id/success", order.CheckOrderCreatedHandler)
	orderGroup.Post("/", order.CreateOrderHandler)

	orderGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	orderGroup.Delete("/", order.DeleteOrdersHandler)

	reviewGroup := v1.Group("/reviews")
	reviewGroup.Post("/", review.CreateReviewHandler)
	reviewGroup.Delete("/", review.BulkDeleteReviewsHandler)

	reviewGroup.Get("/products/:id", review.GetReviewsByProductHandler)
	reviewGroup.Get("/products/:id/overview", review.GetOverviewByProductHandler)

	customerGroup := v1.Group("/customers")
	customerGroup.Get("/", customer.GetCustomersHandler)

	fileGroup := v1.Group("/files")
	fileGroup.Get("/", file.GetFilesHandler)
	fileGroup.Post("/", file.CreateFileHandler)
	fileGroup.Delete("/", file.DeleteFilesHandler)
}
