package routes

import (
	"go-get-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Go Get API is running")
	})

	api := app.Group("/api")

	api.Get("/films", controllers.GetAllFilms)
	api.Get("/films/:id", controllers.GetFilmByID)
	api.Post("/films", controllers.CreateFilm)
	api.Put("/films/:id", controllers.UpdateFilm)
	api.Delete("/films/:id", controllers.DeleteFilm)

	// Endpoint untuk Jadwal
	api.Get("/jadwals", controllers.GetAllJadwals)
	api.Get("/jadwals/:id", controllers.GetJadwalByID)
	app.Get("/jadwals/detail", controllers.GetAllJadwalsWithFilm)
	api.Get("/jadwals/film/:filmId", controllers.GetJadwalsByFilmID)
	api.Post("/jadwals", controllers.CreateJadwal)
	api.Put("/jadwals/:id", controllers.UpdateJadwal)
	api.Delete("/jadwals/:id", controllers.DeleteJadwal)

	// Endpoint untuk Tiket
	api.Get("/tikets", controllers.GetAllTikets)
	api.Get("/tikets/:id", controllers.GetTiketByID)
	api.Post("/tikets", controllers.CreateTiket)
	api.Put("/tikets/:id", controllers.UpdateTiket)
	api.Delete("/tikets/:id", controllers.DeleteTiket)

	// Endpoint untuk Pembayaran
	api.Get("/pembayarans", controllers.GetAllPembayaran)
	api.Get("/pembayarans/:id", controllers.GetPembayaranByID)
	api.Post("/pembayarans", controllers.CreatePembayaran)
	api.Put("/pembayarans/:id", controllers.UpdatePembayaran)
	api.Delete("/pembayarans/:id", controllers.DeletePembayaran)

}
