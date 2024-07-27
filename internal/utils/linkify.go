package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func Linkify(inputText string) string {
	// Define the URL pattern
	urlPattern := regexp.MustCompile(`(?i)(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])|(\bwww\.[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])`)

	// Split the input text by <a> tags
	parts := strings.Split(inputText, `(?i)<a\s+(?:[^>]*?\s+)?href=["'][^"']*["'][^>]*>.*?<\/a>`)

	var linkedText string
	for _, part := range parts {
		if strings.Contains(part, `<a`) {
			// If the part is an <a> tag, leave it unchanged
			linkedText += part
		} else {
			// Otherwise, replace URLs with <a> tags
			linkedText += urlPattern.ReplaceAllStringFunc(part, func(url string) string {
				// Add the protocol if it's missing (for URLs starting with www)
				if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
					url = "http://" + url
				}
				return fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, url, url)
			})
		}
	}
	return linkedText
}
