package main

import (
	"flag"
	"image/png"
	"log"
	"strconv"

	"github.com/a-h/templ"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/qr"
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

	app.Post("/barcode", func(c *fiber.Ctx) error {

		w, h := strconv.FormatInt(barcodeWidth, 10), strconv.FormatInt(barcodeHeight, 10)
		var barcodes []string

		startingNumberStr := c.FormValue("startingNumber")
		barcodeType := c.FormValue("barcodeType")
		numberOfBarcodesStr := c.FormValue("numberOfBarcodes")
		numberOfBarcodes, err := strconv.Atoi(numberOfBarcodesStr)
		if err != nil || numberOfBarcodes < 2 {
			// just do a single barcode
			barcodes = append(barcodes, startingNumberStr)

			cmp := templates.Barcode(barcodes, barcodeType, w, h)

			handler := adaptor.HTTPHandler(templ.Handler(cmp))
			return handler(c)
		}

		prefix, startingNumber, nDigits := splitBarcodePrefix(startingNumberStr)

		// generate 10 barcodes starting from the number
		for i := 0; i < numberOfBarcodes; i++ {
			bc := strconv.Itoa(startingNumber + i)
			for len(bc) < nDigits {
				bc = "0" + bc
			}

			barcodes = append(barcodes, prefix+bc)
		}

		cmp := templates.Barcode(barcodes, barcodeType, w, h)

		handler := adaptor.HTTPHandler(templ.Handler(cmp))
		return handler(c)
	})

	app.Get("/barcode", func(c *fiber.Ctx) error {

		data := c.Query("d")
		barcodeType := c.Query("t")

		var bc barcode.Barcode
		var err error

		switch barcodeType {
		case "Code128":
			bc, err = code128.Encode(data)
		case "Datamatrix":
			bc, err = datamatrix.Encode(data)
		case "QR":
			bc, err = qr.Encode(data, qr.M, qr.Auto)

		}

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		scaledBC, err := barcode.Scale(bc, barcodeWidth, barcodeHeight)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set(fiber.HeaderContentType, "image/png")

		return png.Encode(c.Response().BodyWriter(), scaledBC)
	})

	log.Fatalln(app.Listen(":" + strconv.Itoa(*port)))
}

func splitBarcodePrefix(data string) (prefix string, number int, nDigits int) {

	for i, r := range data {
		if r >= '0' && r <= '9' {
			prefix = data[:i]
			number, _ = strconv.Atoi(data[i:])
			nDigits = len(data[i:])
			break
		}
	}

	return
}
