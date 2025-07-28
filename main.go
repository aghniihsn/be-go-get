package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"go-get-backend/config"
	_ "go-get-backend/docs" // Swagger docs
	"go-get-backend/routes"
)

// @title Cinema Booking API
// @version 1.0
// @description Cinema ticket booking system with payment integration
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host movietix.irc-enter.tech/api
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @schemes http

func main() {
	godotenv.Load()    // load .env
	config.ConnectDB() // koneksi MongoDB

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	// Swagger endpoint
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	// Serve static files for local development
	app.Static("/uploads", "./uploads")

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
