package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/cory-evans/barcode-gen/internal/barcodes"
	"github.com/cory-evans/barcode-gen/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	port := flag.Int("port", 3000, "port to listen on")
	itemsJsonPath := flag.String("itemsJsonPath", "./tmp/items.json", "path to items.json")

	flag.Parse()

	barcodes.FilePath = *itemsJsonPath

	app.Static("/assets", "./assets")

	handlers.Setup(app)

	log.Fatalln(app.Listen(":" + strconv.Itoa(*port)))
}
