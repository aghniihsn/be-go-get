package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"go-get-backend/config"
	"go-get-backend/routes"
)

func main() {
	godotenv.Load()    // load .env
	config.ConnectDB() // koneksi MongoDB

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
