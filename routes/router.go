package routes

import (
	"go-get-backend/config/middleware"
	"go-get-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Go Get API is running")
	})

	api := app.Group("/api")

	// Authentication routes (no auth required)
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/profile", middleware.AuthRequired(), controllers.GetProfile)

	// Public routes (no auth required)
	// Films - public untuk melihat daftar film
	api.Get("/films", controllers.GetAllFilms)
	api.Get("/films/:id", controllers.GetFilmByID)

	// Jadwals - public untuk melihat jadwal
	api.Get("/jadwals", controllers.GetAllJadwals)
	api.Get("/jadwals/:id", controllers.GetJadwalByID)
	app.Get("/jadwals/detail", controllers.GetAllJadwalsWithFilm)
	api.Get("/jadwals/film/:filmId", controllers.GetJadwalsByFilmID)

	// Protected routes - require authentication
	protected := api.Group("", middleware.AuthRequired())

	// User routes (authenticated users only)
	protected.Get("/users/:id", controllers.GetUserByID)
	protected.Put("/users/:id", controllers.UpdateUser)

	// Tiket routes (authenticated users)
	protected.Get("/tikets/user/:user_id", controllers.GetTiketByUserID)
	protected.Post("/tikets", controllers.CreateTiket)
	protected.Put("/tikets/:id", controllers.UpdateTiket)
	protected.Delete("/tikets/:id", controllers.DeleteTiket)

	// Pembayaran routes (authenticated users)
	protected.Get("/pembayarans/:id", controllers.GetPembayaranByID)
	protected.Post("/pembayarans", controllers.CreatePembayaran)
	protected.Put("/pembayarans/:id", controllers.UpdatePembayaran)

	// Admin only routes
	admin := api.Group("", middleware.AdminOnly())

	// Film management (admin only)
	admin.Post("/films", controllers.CreateFilm)
	admin.Put("/films/:id", controllers.UpdateFilm)
	admin.Delete("/films/:id", controllers.DeleteFilm)

	// Jadwal management (admin only)
	admin.Post("/jadwals", controllers.CreateJadwal)
	admin.Put("/jadwals/:id", controllers.UpdateJadwal)
	admin.Delete("/jadwals/:id", controllers.DeleteJadwal)

	// Admin views (admin only)
	admin.Get("/tikets", controllers.GetAllTikets)
	admin.Get("/tikets/:id", controllers.GetTiketByID)
	admin.Get("/pembayarans", controllers.GetAllPembayaran)
	admin.Delete("/pembayarans/:id", controllers.DeletePembayaran)

	// Deprecated route
	api.Post("/users", controllers.CreateUser) // Will return error message
}
