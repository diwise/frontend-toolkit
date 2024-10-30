package assets

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	fetk "github.com/diwise/frontend-toolkit"
)

type AssetHandlerFunc func(fetk.AssetLoader) http.HandlerFunc

type redirect struct {
	path    string
	handler AssetHandlerFunc
}

type ep_conf struct {
	mux *http.ServeMux

	cacheControlHeader string

	redirects []redirect
}

type EndpointOption func(*ep_conf)

func WithImmutableExpiry(maxAge time.Duration) EndpointOption {
	return func(cfg *ep_conf) {
		seconds := int(math.RoundToEven(maxAge.Seconds()))
		cfg.cacheControlHeader = fmt.Sprintf("public,max-age=%d,immutable", seconds)
	}
}

func WithMux(mux *http.ServeMux) EndpointOption {
	return func(cfg *ep_conf) {
		cfg.mux = mux
	}
}

func WithRedirect(from, to string, code int) EndpointOption {
	wildcards := extractWildcards(from)

	return func(cfg *ep_conf) {
		handler := func(assetLoader fetk.AssetLoader) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				redirectPath := to
				for _, wc := range wildcards {
					redirectPath = strings.ReplaceAll(redirectPath, "{"+wc+"}", r.PathValue(wc))
				}

				redirectPath = assetLoader.Load(redirectPath).Path()
				http.Redirect(w, r, redirectPath, code)
			}
		}

		cfg.redirects = append(cfg.redirects, redirect{path: from, handler: handler})
	}
}

func RegisterEndpoints(_ context.Context, assetLoader fetk.AssetLoader, opts ...EndpointOption) error {

	cfg := &ep_conf{
		mux: http.DefaultServeMux,
	}

	for _, applyOption := range opts {
		applyOption(cfg)
	}

	for _, redirect := range cfg.redirects {
		cfg.mux.HandleFunc("GET "+redirect.path, redirect.handler(assetLoader))
	}

	cfg.mux.HandleFunc("GET /assets/{sha}/", func(w http.ResponseWriter, r *http.Request) {
		sha := r.PathValue("sha")

		a, err := assetLoader.LoadFromSha256(sha)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		w.Header().Set("Content-Type", a.ContentType())
		w.Header().Set("Content-Length", fmt.Sprintf("%d", a.ContentLength()))

		if cfg.cacheControlHeader != "" {
			w.Header().Set("Cache-Control", cfg.cacheControlHeader)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(a.Body())
	})

	return nil
}

func extractWildcards(path string) []string {
	wildcards := make([]string, 0, 5)
	curlyIdx := strings.Index(path, "{")

	for curlyIdx >= 0 {
		path = path[curlyIdx+1:]
		curlyIdx = strings.Index(path, "}")
		if curlyIdx >= 0 {
			wildcards = append(wildcards, path[0:curlyIdx])
			path = path[curlyIdx:]

			curlyIdx = strings.Index(path, "{")
		}
	}

	return wildcards
}
