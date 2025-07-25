package controllers

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// SwaggerHandler returns Fiber handler for Swagger UI
func SwaggerHandler() fiber.Handler {
	return fiberSwagger.WrapHandler
}
