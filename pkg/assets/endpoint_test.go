package assets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	frontendtoolkit "github.com/diwise/frontend-toolkit"
	"github.com/diwise/frontend-toolkit/mock"

	"github.com/matryer/is"
)

func TestWC(t *testing.T) {
	is := is.New(t)
	wc := extractWildcards("/assets/{sha}/{filename}")

	is.Equal(len(wc), 2)
	is.Equal(wc[0], "sha")
	is.Equal(wc[1], "filename")
}

func TestRegisterEndpoints(t *testing.T) {
	is := is.New(t)
	mux := http.NewServeMux()
	asset := &mock.AssetMock{
		BodyFunc:          func() []byte { return make([]byte, 10) },
		ContentLengthFunc: func() int { return 10 },
		ContentTypeFunc:   func() string { return "application/png" },
		PathFunc: func() string {
			return "/assets/1337/file.png"
		},
	}

	loader := &mock.LoaderMock{
		LoadFunc:           func(name string) frontendtoolkit.Asset { return asset },
		LoadFromSha256Func: func(string) (frontendtoolkit.Asset, error) { return asset, nil },
	}

	err := RegisterEndpoints(context.Background(), loader,
		WithMux(mux),
		WithRedirect("/favicon.ico", "/icons/favicon.ico", http.StatusMovedPermanently),
		WithRedirect("/assets/1337/images/{img}", "/images/leaflet-{img}", http.StatusFound),
	)
	is.NoErr(err)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/favicon.ico")
	is.NoErr(err)
	is.Equal(resp.StatusCode, http.StatusOK)

	resp, err = http.Get(ts.URL + "/assets/1337/images/marker.png")
	is.NoErr(err)
	is.Equal(resp.StatusCode, http.StatusOK)
}
