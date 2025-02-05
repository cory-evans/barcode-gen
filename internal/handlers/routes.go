package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Get("/", Index)
	app.Post("/generate", GenerateOutputHTML)
	app.Post("/load", LoadOutputHTML)
	app.Get("/barcode", GenerateBarcodeImage)

	savedItems := app.Group("/saved")

	savedItems.Get("/search", SearchSavedItems)
	savedItems.Get("/list", ListSavedItems)
	savedItems.Delete("/delete", DeleteSavedItem)
}
