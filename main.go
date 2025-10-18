package main

import (
	"app/cmd/server"
	"app/internal/database"
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
	database.InitDB("./db.sqlite")
	server.Serve()
}
