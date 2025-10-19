package router

import (
	"app/internal/modules/auth"
	"app/internal/modules/file"
	"app/internal/modules/user"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	v1 := app.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authGroup.Post("/signin", auth.SignInHandler)
	authGroup.Post("/signup", auth.SignUpHandler)

	userGroup := v1.Group("/users")
	userGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	userGroup.Get("/", user.GetUsersHandler)
	userGroup.Get("/:id", user.GetUserHandler)
	userGroup.Post("/", user.CreateUserHandler)
	userGroup.Delete("/", user.DeleteUsersHandler)

	fileGroup := v1.Group("/files")
	fileGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("jwt")},
	}))
	fileGroup.Post("/upload", file.UploadFilesHanlder)
	fileGroup.Delete("/delete", file.DeleteFilesHandler)
}
