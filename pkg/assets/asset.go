package assets

import (
	"errors"
	"path/filepath"
)

var ErrNotFound = errors.New("not found")

type asset struct {
	path        string
	sha256      string
	contentType string
	body        []byte
}

func contentTypeFromFileName(filename string) string {

	extension := filepath.Ext(filename)

	contentType, ok := map[string]string{
		".css":   "text/css; charset=utf-8",
		".ico":   "image/x-icon",
		".js":    "text/javascript; charset=utf-8",
		".png":   "image/png",
		".svg":   "image/svg+xml",
		".webp":  "image/webp",
		".woff2": "font/woff2",
	}[extension]

	if !ok {
		return "text/plain"
	}

	return contentType
}

func (a asset) Body() []byte {
	return a.body
}

func (a asset) ContentLength() int {
	return len(a.body)
}

func (a asset) ContentType() string {
	return a.contentType
}

func (a asset) Path() string {
	return a.path
}

func (a asset) SHA256() string {
	return a.sha256
}
