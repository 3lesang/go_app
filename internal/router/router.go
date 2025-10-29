package router

import (
	"app/internal/modules/auth"
	"app/internal/modules/category"
	"app/internal/modules/collection"
	"app/internal/modules/product"
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
	productGroup.Get("/", product.GetProductsHandler)
	productGroup.Get("/:id", product.GetProductHandler)
	productGroup.Post("/", product.CreateProductHandler)
	productGroup.Put("/:id", product.UpdateProductHandler)
	productGroup.Delete("/", product.DeleteProductsHandler)

	categoryGroup := v1.Group("/categories")
	categoryGroup.Get("/", category.GetCategoriesHandler)
	categoryGroup.Get("/:id", category.GetCategoryHandler)
	categoryGroup.Post("/", category.CreateCategoryHandler)
	categoryGroup.Put("/:id", category.UpdateCategoryHandler)
	categoryGroup.Delete("/", category.DeleteCategoriesHandler)

	collectionGroup := v1.Group("/collections")
	collectionGroup.Get("/", collection.GetCollectionsHandler)
	collectionGroup.Get("/:id", collection.GetCollectionHandler)
	collectionGroup.Post("/", collection.CreateCollectionHandler)
	collectionGroup.Put("/:id", collection.UpdateCollectionHandler)
	collectionGroup.Delete("/", collection.DeleteCollectionsHandler)
}
