package handlers

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/cory-evans/barcode-gen/internal/barcodes"
	"github.com/cory-evans/barcode-gen/internal/components"
	"github.com/cory-evans/barcode-gen/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func SearchSavedItems(c *fiber.Ctx) error {
	query := strings.ToLower(c.Query("search"))

	if query == "" {
		return ListSavedItems(c)
	}

	items := barcodes.GetItems()
	filtered := []models.GenerateInput{}

	for _, v := range items {
		if strings.Contains(strings.ToLower(v.SaveName), query) {
			filtered = append(filtered, v)
		}
	}

	cmp := components.SavedItemList(filtered)

	h := adaptor.HTTPHandler(templ.Handler(cmp))
	return h(c)
}

func ListSavedItems(c *fiber.Ctx) error {
	items := barcodes.GetItems()
	cmp := components.SavedItemList(items)

	h := adaptor.HTTPHandler(templ.Handler(cmp))
	return h(c)
}

func DeleteSavedItem(c *fiber.Ctx) error {
	name := c.Query("save_name")
	barcodes.DeleteItem(name)

	c.Response().Header.Set("HX-Trigger", "refresh_saved")

	return ListSavedItems(c)
}
