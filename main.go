package main

import (
	"app/cmd/server"
	"app/internal/db"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath  /api/v1
func main() {
	db.InitDB("./sqlite.db")
	server.Serve()
}
