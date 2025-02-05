package handlers

import (
	"fmt"
	"image/png"
	"strings"

	"github.com/a-h/templ"
	"github.com/cory-evans/barcode-gen/internal/barcodes"
	"github.com/cory-evans/barcode-gen/internal/components"
	"github.com/cory-evans/barcode-gen/internal/models"
	"github.com/cory-evans/barcode-gen/internal/util"
	"github.com/cory-evans/barcode-gen/pkg/array"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func getBarcodeData(input models.GenerateInput) []models.BarcodeData {
	// generate barcode data
	var barcodes []string
	if input.GenerateType == "Sequence" {

		for i := 0; i < input.Number; i++ {
			d := fmt.Sprintf("%s%d%s", input.Prefix, input.Start+i, input.Suffix)

			barcodes = append(barcodes, d)
		}

	} else if input.GenerateType == "OnePerLine" {

		barcodes = strings.Split(strings.TrimSpace(input.BarcodeData), "\n")
		barcodes = array.Filter(barcodes, func(b string) bool {
			return b != ""
		})

		for i := 0; i < len(barcodes); i++ {
			barcodes[i] = fmt.Sprintf("%s%s%s", input.Prefix, barcodes[i], input.Suffix)
		}
	}

	if input.Width == 0 {
		input.Width = 200
	}

	if input.Height == 0 {
		input.Height = 50
	}

	bd := array.Map(barcodes, func(b string) models.BarcodeData {
		return models.BarcodeData{
			Type:   input.BarcodeType,
			Data:   b,
			Width:  input.WidthAsStr(),
			Height: input.HeightAsStr(),
		}
	})

	return bd
}

func GenerateOutputHTML(c *fiber.Ctx) error {

	input := models.NewGenerateInputFromForm(c)

	if input.Action == "save" {
		barcodes.SaveItem(input)
		c.Response().Header.Set("HX-Trigger", "refresh_saved")
	}

	cmp := components.BarcodeOutput(getBarcodeData(input))

	h := adaptor.HTTPHandler(templ.Handler(cmp))

	return h(c)
}

func LoadOutputHTML(c *fiber.Ctx) error {
	saveName := c.Query("save_name")
	input := *barcodes.GetItem(saveName)
	cmp := components.BarcodeOutputWithForm(components.NewBarcodeFormProps(input), getBarcodeData(input))
	h := adaptor.HTTPHandler(templ.Handler(cmp))
	return h(c)
}

func GenerateBarcodeImage(c *fiber.Ctx) error {

	data := c.Query("d")
	barcodeType := c.Query("t")
	w, h := util.FormValueAsInt(c, "w"), util.FormValueAsInt(c, "h")

	bc, err := barcodes.Generate(barcodeType, data, w, h)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set(fiber.HeaderContentType, "image/png")

	return png.Encode(c.Response().BodyWriter(), bc)
}
