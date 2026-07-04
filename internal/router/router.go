package router

import (
	"github.com/Zigl3ur/api-portfolio/internal/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func SetupRoutes(app *fiber.App, musicHandler *handlers.MusicHandler, messageHandler *handlers.MessageHandler) {

	musicGroup := app.Group("/music")
	musicGroup.Get("/currently-listening", musicHandler.CurrentlyListening)
	musicGroup.Get("/top-albums", musicHandler.TopAlbums)

	app.Post("/message", handlers.MessageLimiter, messageHandler.Handler)

	app.Get("/videos*", static.New("./videos"))
}
