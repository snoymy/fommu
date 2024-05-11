package utils

import (
	"net/http"
	"strings"
)

func GetClientPlatform(r *http.Request) (string, string) {
    userAgent := r.UserAgent()

    var device string
    if strings.Contains(userAgent, "Mobi") {
        device = "Mobile"
    } else {
        device = "Desktop"
    }

    var os string
    if strings.Contains(userAgent, "Android") {
        os = "Android"
    } else if strings.Contains(userAgent, "iPhone") {
        os = "iPhone"
    } else if strings.Contains(userAgent, "Windows") {
        os = "Windows"
    } else if strings.Contains(userAgent, "Macintosh") {
        os = "Macintosh"
    } else if strings.Contains(userAgent, "Linux") {
        os = "Linux"
    } else {
        os = "Unknown"
    }
    
    return device, os
}
