package router

import (
	"app/internal/modules/auth"
	"app/internal/modules/category"
	"app/internal/modules/collection"
	"app/internal/modules/customer"
	"app/internal/modules/discount"
	"app/internal/modules/file"
	"app/internal/modules/hotspot"
	"app/internal/modules/menu"
	"app/internal/modules/order"
	"app/internal/modules/page"
	"app/internal/modules/post"
	"app/internal/modules/product"
	"app/internal/modules/review"
	"app/internal/modules/search"
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
	productGroup.Get("/categories/:id", product.GetProductByCategoryHandler)

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

	customerGroup.Post("/register", customer.RegisterCustomerHandler)
	customerGroup.Post("/login", customer.CustomerLoginHandler)

	customerGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	customerGroup.Get("/", customer.GetCustomersHandler)
	customerGroup.Get("/me", customer.GetMeHandler)
	customerGroup.Post("/me", customer.UpdateMeHandler)
	customerGroup.Post("/", customer.CreateCustomerHandler)
	customerGroup.Delete("/", customer.BulkDeleteCustomersHandler)

	pageGroup := v1.Group("/pages")
	pageGroup.Get("/", page.GetPagesHandler)
	pageGroup.Get("/:id", page.GetPageHandler)
	pageGroup.Get("/slug/:id", page.GetPageBySlugHandler)
	pageGroup.Post("/", page.CreatePageHandler)
	pageGroup.Put("/:id", page.UpdatePageHandler)
	pageGroup.Delete("/", page.BulkDeletePagesHandler)

	postGroup := v1.Group("/posts")
	postGroup.Get("/", post.GetPostsHandler)
	postGroup.Get("/:id", post.GetPostHandler)
	postGroup.Get("/slug/:id", post.GetPostBySlugHandler)
	postGroup.Post("/", post.CreatePostHandler)
	postGroup.Put("/:id", post.UpdatePostHandler)
	postGroup.Delete("/", post.BulkDeletePostsHandler)

	menuGroup := v1.Group("/menus")
	menuGroup.Get("/", menu.GetMenusHandler)
	menuGroup.Get("/:id", menu.GetMenuHandler)
	menuGroup.Get("/position/:id", menu.GetMenuByPositionHandler)
	menuGroup.Post("/", menu.CreateMenuHandler)
	menuGroup.Put("/:id", menu.UpdateMenuHandler)
	menuGroup.Delete("/", menu.BulkDeleteMenusHandler)

	fileGroup := v1.Group("/files")
	fileGroup.Get("/", file.GetFilesHandler)
	fileGroup.Post("/", file.CreateFileHandler)
	fileGroup.Delete("/", file.DeleteFilesHandler)

	discountGroup := v1.Group("/discounts")
	discountGroup.Get("/public", discount.GetValidDiscountsHandler)

	discountGroup.Get("/", discount.GetDiscountsHandler)
	discountGroup.Get("/:id", discount.GetDiscountHandler)
	discountGroup.Put("/:id", discount.UpdateDiscountHandler)
	discountGroup.Post("/", discount.CreateDiscountHandler)
	discountGroup.Delete("/", discount.BulkDeleteDiscountsHandler)

	discountGroup.Get("/:discount_id/customers/:customer_id/usage", discount.GetCustomerUsageHandler)
	discountGroup.Post("/usage", discount.UpsertCustomerUsageHandler)

	discountGroup.Post("/:id/targets", discount.CreateDiscountTargetHandler)
	discountGroup.Post("/:id/effects", discount.CreateDiscountEffectHandler)
	discountGroup.Post("/:id/conditions", discount.CreateDiscountConditionHandler)
	discountGroup.Put("/:id/effects/:effectID", discount.UpdateDiscountEffectHandler)
	discountGroup.Put("/:id/conditions/:conditionID", discount.UpdateDiscountConditionHandler)

	hotspotGroup := v1.Group("/hotspots")
	hotspotGroup.Get("/", hotspot.GetHotspotsHandler)
	hotspotGroup.Get("/:id", hotspot.GetHotspotHandler)
	hotspotGroup.Get("/products/:id", hotspot.GetHotspotByProductHandler)
	hotspotGroup.Post("/", hotspot.CreateHotspotHandler)
	hotspotGroup.Put("/:id", hotspot.UpdateHotspotHandler)
	hotspotGroup.Delete("/", hotspot.DeleteHotspotsHandler)

	v1.Get("/search", search.SearchProductsHandler)
}
