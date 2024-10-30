package csp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestCSPWithStrictHeader(t *testing.T) {
	is := is.New(t)
	wr := httptest.NewRecorder()

	uut := NewContentSecurityPolicy()
	h := uut(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	h.ServeHTTP(wr, httptest.NewRequest(http.MethodGet, "/", nil))

	value, exists := wr.Header()[ContentSecurityPolicy]
	is.True(exists)
	is.Equal(len(value), 1) // should consist of a single value

	// remove the random nonce from the value before matching
	valueWithoutNonce := value[0][0:17] + value[0][50:]

	is.Equal(valueWithoutNonce, "script-src 'nonce'; object-src 'none'; base-uri 'none';")
}

func TestCSPWithStrictDynamicHeader(t *testing.T) {
	is := is.New(t)
	wr := httptest.NewRecorder()

	uut := NewContentSecurityPolicy(StrictDynamic())
	h := uut(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	h.ServeHTTP(wr, httptest.NewRequest(http.MethodGet, "/", nil))

	value, exists := wr.Header()[ContentSecurityPolicy]
	is.True(exists)
	is.Equal(len(value), 1) // should consist of a single value

	// remove the random nonce from the value before matching
	valueWithoutNonce := value[0][0:17] + value[0][50:]

	is.Equal(valueWithoutNonce, "script-src 'nonce' 'strict-dynamic'; object-src 'none'; base-uri 'none';")
}

func TestCSPWithStrictDynamicAndReportToHeader(t *testing.T) {
	is := is.New(t)
	wr := httptest.NewRecorder()

	uut := NewContentSecurityPolicy(StrictDynamic(), ReportTo("https://violations.r.us"))
	h := uut(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	h.ServeHTTP(wr, httptest.NewRequest(http.MethodGet, "/", nil))

	cspValue, exists := wr.Header()[ContentSecurityPolicy]
	is.True(exists)            // csp header should have been set
	is.Equal(len(cspValue), 1) // should consist of a single value

	reValue, exists := wr.Header()[ReportingEndpoints]
	is.True(exists)           // reporting endpoints header should have been set
	is.Equal(len(reValue), 1) // should consist of a single value

	// remove the random nonce from the value before matching
	valueWithoutNonce := cspValue[0][0:17] + cspValue[0][50:]
	is.Equal(valueWithoutNonce, "script-src 'nonce' 'strict-dynamic'; object-src 'none'; base-uri 'none'; report-to csp-endpoint")
	is.Equal(reValue[0], "csp-endpoint=\"https://violations.r.us\"")
}

func TestCSPWithReportOnlyHeader(t *testing.T) {
	is := is.New(t)
	wr := httptest.NewRecorder()

	uut := NewContentSecurityPolicy(ReportOnly())
	h := uut(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	h.ServeHTTP(wr, httptest.NewRequest(http.MethodGet, "/", nil))

	_, exists := wr.Header()[ContentSecurityPolicyReportOnly]
	is.True(exists)
}

func TestCSPPassesNonceInContext(t *testing.T) {
	is := is.New(t)
	assertionDone := false

	uut := NewContentSecurityPolicy()
	h := uut(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.True(HasNonce(r.Context()))
		is.Equal(len(Nonce(r.Context())), 32) // nonce should be 128 bit (32 bytes hex)
		assertionDone = true
	}))

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	is.True(assertionDone) // csp middleware did not call next handler
}
