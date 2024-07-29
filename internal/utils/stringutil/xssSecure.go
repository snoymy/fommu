package stringutil

import (
	"html"

	"github.com/microcosm-cc/bluemonday"
)

func XSSSecure(input string) string {
    encodedInput := html.EscapeString(input)

    p := bluemonday.UGCPolicy()

    return p.Sanitize(encodedInput)
}
