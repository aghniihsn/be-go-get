package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"go-get-backend/config"
	"go-get-backend/routes"
)

func main() {
	fmt.Println("Starting application...")
	godotenv.Load() // load .env
	fmt.Println("Environment loaded")
	config.ConnectDB() // koneksi MongoDB
	fmt.Println("Database connected")

	app := fiber.New()
	fmt.Println("Fiber app created")

	app.Use(cors.New())
	app.Use(logger.New())
	fmt.Println("Middlewares configured")

	routes.SetupRoutes(app)
	fmt.Println("Routes configured")

	fmt.Println("Starting server on port 3000...")
	err := app.Listen(":3000")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
