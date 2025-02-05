package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cory-evans/barcode-gen/internal/util"
	"github.com/gofiber/fiber/v2"
)

var GenerationTypes = []string{"Sequence", "OnePerLine"}

type GenerateInput struct {
	GenerateType string `json:"generate_type"`
	BarcodeType  string `json:"barcode_type"`
	Prefix       string `json:"prefix"`
	Suffix       string `json:"suffix"`
	BarcodeData  string `json:"barcode_data"`
	Start        int    `json:"start"`
	Number       int    `json:"number"`

	Width  int `json:"width"`
	Height int `json:"height"`

	SaveName string `json:"save_name"`

	Action string `json:"-"`
}

func NewGenerateInputFromForm(c *fiber.Ctx) GenerateInput {
	return GenerateInput{
		GenerateType: c.FormValue("generate_type"),
		BarcodeType:  c.FormValue("barcode_type"),
		Prefix:       c.FormValue("prefix"),
		Suffix:       c.FormValue("suffix"),
		BarcodeData:  c.FormValue("barcode_data"),
		Start:        util.FormValueAsInt(c, "start"),
		Number:       util.FormValueAsInt(c, "number"),

		Width:  util.FormValueAsInt(c, "width"),
		Height: util.FormValueAsInt(c, "height"),

		SaveName: c.FormValue("save_name"),

		Action: c.FormValue("action"),
	}
}

func NewBlankGenerateInput() GenerateInput {
	start, _ := strconv.Atoi(time.Now().Format("20060102") + "01")
	return GenerateInput{
		GenerateType: GenerationTypes[0],
		BarcodeType:  "code128",
		Prefix:       "",
		Suffix:       "",
		BarcodeData:  "",
		Start:        start,
		Number:       1,

		Width:  200,
		Height: 50,

		SaveName: "",
	}
}

func (i *GenerateInput) StartAsStr() string {
	return fmt.Sprintf("%d", i.Start)
}

func (i *GenerateInput) NumberAsStr() string {
	return fmt.Sprintf("%d", i.Number)
}

func (i *GenerateInput) WidthAsStr() string {
	return fmt.Sprintf("%d", i.Width)
}

func (i *GenerateInput) HeightAsStr() string {
	return fmt.Sprintf("%d", i.Height)
}

type BarcodeData struct {
	Type string `json:"type"`
	Data string `json:"data"`

	Width  string `json:"width"`
	Height string `json:"height"`
}

func (bc *BarcodeData) URLParams() string {
	return "?d=" + bc.Data + "&t=" + bc.Type + "&w=" + bc.Width + "&h=" + bc.Height
}
