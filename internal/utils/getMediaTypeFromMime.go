package utils

func GetMediaTypeFromMime(mimeType string) string {
    mimeToMediaType := map[string]string{
        "application/atom+xml":        "document",
        "application/ecmascript":      "document",
        "application/EDI-X12":         "document",
        "application/EDIFACT":         "document",
        "application/json":            "document",
        "application/javascript":      "document",
        "application/octet-stream":    "document",
        "application/ogg":             "audio",
        "application/pdf":             "document",
        "application/xhtml+xml":       "document",
        "application/xml-dtd":         "document",
        "application/zip":             "document",
        "audio/midi":                  "audio",
        "audio/mpeg":                  "audio",
        "audio/ogg":                   "audio",
        "audio/x-wav":                 "audio",
        "image/gif":                   "image",
        "image/jpeg":                  "image",
        "image/png":                   "image",
        "image/svg+xml":               "image",
        "image/tiff":                  "image",
        "image/vnd.microsoft.icon":    "image",
        "text/calendar":               "document",
        "text/css":                    "document",
        "text/csv":                    "document",
        "text/html":                   "document",
        "text/javascript":             "document",
        "text/plain":                  "document",
        "text/xml":                    "document",
        "video/3gpp":                  "video",
        "video/mp4":                   "video",
        "video/mpeg":                  "video",
        "video/quicktime":             "video",
        "video/webm":                  "video",
        "video/x-flv":                 "video",
    }
    if mediaType, ok := mimeToMediaType[mimeType]; ok {
        return mediaType
    } else {
        return "unknown"
    }
}
