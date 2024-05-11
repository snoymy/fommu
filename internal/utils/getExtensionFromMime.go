package utils

func GetExtensionFromMIME(mimeType string) string {
    mimeToExtension := map[string]string{
        "application/atom+xml":           ".atom",
        "application/ecmascript":         ".ecma",
        "application/EDI-X12":            ".edi",
        "application/EDIFACT":            ".edi",
        "application/json":               ".json",
        "application/javascript":         ".js",
        "application/octet-stream":       ".bin",
        "application/ogg":                ".ogg",
        "application/pdf":                ".pdf",
        "application/xhtml+xml":          ".xhtml",
        "application/xml-dtd":            ".dtd",
        "application/zip":                ".zip",
        "audio/midi":                     ".midi",
        "audio/mpeg":                     ".mp3",
        "audio/ogg":                      ".ogg",
        "audio/x-wav":                    ".wav",
        "image/gif":                      ".gif",
        "image/jpeg":                     ".jpg",
        "image/png":                      ".png",
        "image/svg+xml":                  ".svg",
        "image/tiff":                     ".tiff",
        "image/vnd.microsoft.icon":       ".ico",
        "text/calendar":                  ".ics",
        "text/css":                       ".css",
        "text/csv":                       ".csv",
        "text/html":                      ".html",
        "text/javascript":                ".js",
        "text/plain":                     ".txt",
        "text/xml":                       ".xml",
        "video/3gpp":                     ".3gp",
        "video/mp4":                      ".mp4",
        "video/mpeg":                     ".mpeg",
        "video/quicktime":                ".mov",
        "video/webm":                     ".webm",
        "video/x-flv":                    ".flv",
    }
    if extension, ok := mimeToExtension[mimeType]; ok {
        return extension
    } else {
        return ""
    }
}
