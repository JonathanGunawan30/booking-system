package main

import "user-service/cmd"

// @title User Service API
// @version 1.0
// @description Microservice for user management and authentication.
// @host localhost:8001
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
func main() {
	cmd.Run()
}
