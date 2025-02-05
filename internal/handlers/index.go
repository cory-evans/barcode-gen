package handlers

import (
	"github.com/a-h/templ"
	"github.com/cory-evans/barcode-gen/internal/templates"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Index(c *fiber.Ctx) error {
	cmp := templates.Home()

	handler := adaptor.HTTPHandler(templ.Handler(cmp))
	return handler(c)
}
