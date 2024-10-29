package assets

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	tk "github.com/diwise/frontend-toolkit"
)

type impl struct {
	basePath string

	path2sha map[string]string
	assets   map[string]tk.Asset

	exclude []func(string, string) bool

	logger *slog.Logger

	mu sync.Mutex
}

type LoaderOptionFunc func(*impl)

func BasePath(base string) LoaderOptionFunc {
	return func(loader *impl) {
		loader.basePath = base
	}
}

func Exclude(match string) LoaderOptionFunc {
	return func(loader *impl) {
		loader.exclude = append(loader.exclude, func(path, name string) bool {
			return strings.Contains(path, match) || strings.Contains(name, match)
		})
	}
}

func Logger(logger *slog.Logger) LoaderOptionFunc {
	return func(loader *impl) {
		loader.logger = logger
	}
}

func NewLoader(ctx context.Context, opts ...LoaderOptionFunc) (tk.AssetLoader, error) {
	loader := &impl{
		basePath: ".",
		path2sha: map[string]string{},
		assets:   map[string]tk.Asset{},
		exclude:  []func(string, string) bool{},
	}

	for _, opt := range append(opts, Exclude(".DS_Store")) {
		opt(loader)
	}

	filepath.Walk(loader.basePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			for _, shouldBeExcluded := range loader.exclude {
				if shouldBeExcluded(path, info.Name()) {
					return nil
				}
			}

			loader.Load(
				strings.TrimPrefix(path, loader.basePath),
			)

			return nil
		})

	return loader, nil
}

func (l *impl) Load(path string) tk.Asset {
	l.mu.Lock()
	defer l.mu.Unlock()

	if sha, ok := l.path2sha[path]; ok {
		return l.assets[sha]
	}

	fullPath := l.basePath + path

	assetFile, err := os.Open(fullPath)
	if err != nil {
		panic("failed to open asset " + fullPath + " (" + err.Error() + ")")
	}
	defer assetFile.Close()

	assetContents, err := io.ReadAll(assetFile)
	if err != nil {
		panic("failed to read from asset file " + fullPath + " (" + err.Error() + ")")
	}

	sha256 := fmt.Sprintf("%x", sha256.Sum256(assetContents))
	l.path2sha[path] = sha256

	pathTokens := strings.Split(path, "/")
	assetFileName := pathTokens[len(pathTokens)-1]

	a := &asset{
		path:        fmt.Sprintf("/assets/%s/%s", sha256, assetFileName),
		contentType: contentTypeFromFileName(path),
		body:        assetContents,
		sha256:      sha256,
	}

	l.assets[sha256] = a

	if l.logger != nil {
		l.logger.Debug("loaded " + path + " as " + a.path)
	}

	return a
}

func (l *impl) LoadFromSha256(sha string) (tk.Asset, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if a, ok := l.assets[sha]; ok {
		return a, nil
	}

	return nil, ErrNotFound
}
