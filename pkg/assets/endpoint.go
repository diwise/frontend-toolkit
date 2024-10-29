package assets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	fetk "github.com/diwise/frontend-toolkit"
)

type AssetHandlerFunc func(fetk.AssetLoader) http.HandlerFunc

type redirect struct {
	path    string
	handler AssetHandlerFunc
}

type ep_conf struct {
	mux *http.ServeMux

	redirects []redirect
}

type EndpointOption func(*ep_conf)

func WithMux(mux *http.ServeMux) EndpointOption {
	return func(cfg *ep_conf) {
		cfg.mux = mux
	}
}

// /favicon.ico -> /icons/favicon.ico http.StatusMovedPermanently
// 	WithRedirect("/favicon.ico", "/icons/favicon.ico", http.StatusMovedPermanently)
// /assets/<sha>/images/{img} -> /images/leaflet-"+image http.StatusFound
// 	WithRedirect()

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
