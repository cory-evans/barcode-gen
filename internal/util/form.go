package util

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FormValueAsInt(c *fiber.Ctx, key string) int {
	v := c.FormValue(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}
