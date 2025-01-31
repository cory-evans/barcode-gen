package main

import (
	"flag"
	"image/png"
	"log"
	"strconv"

	"github.com/a-h/templ"
	"github.com/cory-evans/barcode-gen/internal/barcodes"
	"github.com/cory-evans/barcode-gen/internal/components"
	"github.com/cory-evans/barcode-gen/internal/templates"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

const barcodeWidth = 300
const barcodeHeight = 150

func main() {
	app := fiber.New()

	port := flag.Int("port", 3000, "port to listen on")
	flag.Parse()

	app.Static("/assets", "./assets")

	app.Get("/", func(c *fiber.Ctx) error {
		cmp := templates.Home()

		handler := adaptor.HTTPHandler(templ.Handler(cmp))
		return handler(c)
	})
	app.Get("/form", func(c *fiber.Ctx) error {

		idstr := c.QueryInt("id", 0)
		cmp := components.BarcodeForm(idstr)

		handler := adaptor.HTTPHandler(templ.Handler(cmp))
		return handler(c)
	})

	app.Post("/barcode", func(c *fiber.Ctx) error {

		w, h := strconv.FormatInt(barcodeWidth, 10), strconv.FormatInt(barcodeHeight, 10)
		var barcodeData []string

		startingNumberStr := c.FormValue("startingNumber")
		barcodeType := c.FormValue("barcodeType")
		numberOfBarcodesStr := c.FormValue("numberOfBarcodes")
		numberOfBarcodes, err := strconv.Atoi(numberOfBarcodesStr)
		if err != nil || numberOfBarcodes < 2 {
			// just do a single barcode
			barcodeData = append(barcodeData, startingNumberStr)

			cmp := templates.Barcode(barcodeData, barcodeType, w, h)

			handler := adaptor.HTTPHandler(templ.Handler(cmp))
			return handler(c)
		}

		prefix, startingNumber, nDigits := barcodes.SplitBarcodePrefix(startingNumberStr)

		// generate 10 barcodes starting from the number
		for i := 0; i < numberOfBarcodes; i++ {
			bc := strconv.Itoa(startingNumber + i)
			for len(bc) < nDigits {
				bc = "0" + bc
			}

			barcodeData = append(barcodeData, prefix+bc)
		}

		cmp := templates.Barcode(barcodeData, barcodeType, w, h)

		handler := adaptor.HTTPHandler(templ.Handler(cmp))
		return handler(c)
	})

	app.Get("/barcode", func(c *fiber.Ctx) error {

		data := c.Query("d")
		barcodeType := c.Query("t")

		bc, err := barcodes.Generate(barcodeType, data, barcodeWidth, barcodeHeight)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set(fiber.HeaderContentType, "image/png")

		return png.Encode(c.Response().BodyWriter(), bc)
	})

	log.Fatalln(app.Listen(":" + strconv.Itoa(*port)))
}
