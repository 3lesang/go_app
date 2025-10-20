package main

import (
	"app/cmd/server"
	database "app/internal/database/postgres"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @basePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	database.Init()
	defer database.Pool.Close()
	server.Serve()
}
