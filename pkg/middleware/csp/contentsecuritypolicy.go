package csp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

const (
	ContentSecurityPolicy           string = "Content-Security-Policy"
	ContentSecurityPolicyReportOnly string = "Content-Security-Policy-Report-Only"
	ReportingEndpoints              string = "Reporting-Endpoints"
)

type nonceKeyType string

const nonceKey nonceKeyType = "nonce-key"

func HasNonce(ctx context.Context) bool {
	return ctx.Value(nonceKey) != nil
}

func Nonce(ctx context.Context) string {
	if v := ctx.Value(nonceKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type cspcfg struct {
	reportOnly         bool
	reportViolationsTo string
	strictDynamic      bool
}

type cspOpt func(*cspcfg)

func ReportOnly() cspOpt {
	return func(cfg *cspcfg) {
		cfg.reportOnly = true
	}
}

func ReportTo(url string) cspOpt {
	return func(cfg *cspcfg) {
		cfg.reportViolationsTo = fmt.Sprintf("csp-endpoint=\"%s\"", url)
	}
}

func StrictDynamic() cspOpt {
	return func(cfg *cspcfg) {
		cfg.strictDynamic = true
	}
}

func NewContentSecurityPolicy(opts ...cspOpt) func(http.Handler) http.Handler {
	cfg := &cspcfg{}
	for _, applyOption := range opts {
		applyOption(cfg)
	}

	cspHeader := map[bool]string{false: ContentSecurityPolicy, true: ContentSecurityPolicyReportOnly}[cfg.reportOnly]

	var sb strings.Builder
	sb.WriteString("script-src 'nonce-%s'")
	if cfg.strictDynamic {
		sb.WriteString(" 'strict-dynamic'")
	}
	sb.WriteString("; object-src 'none'; base-uri 'none';")

	if cfg.reportViolationsTo != "" {
		sb.WriteString(" report-to csp-endpoint")
	}

	policyFmt := sb.String()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nonceBytes := make([]byte, 16)
			if _, err := rand.Read(nonceBytes); err == nil {
				if cfg.reportViolationsTo != "" {
					w.Header().Set(ReportingEndpoints, cfg.reportViolationsTo)
				}
				nonce := hex.EncodeToString(nonceBytes)
				w.Header().Set(cspHeader, fmt.Sprintf(policyFmt, nonce))
				r = r.WithContext(context.WithValue(r.Context(), nonceKey, nonce))
			}

			next.ServeHTTP(w, r)
		})
	}
}
