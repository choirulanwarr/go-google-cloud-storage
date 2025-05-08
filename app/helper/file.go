package helper

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"
)

var allowedMIMETypes = map[string]bool{
	"image/jpeg":               true,
	"image/png":                true,
	"image/jpg":                true,
	"application/pdf":          true,
	"application/zip":          true,
	"application/octet-stream": true,
}

func IsAllowedFileType(apiCallID string, file multipart.File) bool {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false
	}
	_, _ = file.Seek(0, io.SeekStart)

	mimeType := http.DetectContentType(buffer)
	LogInfo(apiCallID, "Content-Type: "+mimeType)
	return allowedMIMETypes[mimeType]
}

func DefaultMIME(mime string) string {
	if mime == "" {
		return "application/octet-stream"
	}
	return mime
}

func GenerateUniqueFilename() string {
	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d%02d", now.Unix(), now.Hour(), now.Minute(), now.Second())
	randomSuffix := rand.Intn(10_000_000) + 1
	return fmt.Sprintf("%s%d", timestamp, randomSuffix)
}

func GeneratePublicURL(bucketName, filename string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
}

func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
