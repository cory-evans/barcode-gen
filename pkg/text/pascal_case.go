package text

import (
	"strings"
	"unicode"
)

func PascalCase(s string) string {

	var words []string

	i := 0
	for x := s; x != ""; x = x[i:] {
		i = strings.IndexFunc(x[1:], unicode.IsUpper) + 1
		if i <= 0 {
			i = len(x)
		}

		words = append(words, x[:i])
	}

	return strings.Join(words, " ")
}
