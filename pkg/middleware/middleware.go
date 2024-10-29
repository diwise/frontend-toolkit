package middleware

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

func StrictTransportSecurity(maxAge time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		const TwoYears int = 2 * 365 * 24 * 60 * 60

		seconds := int(math.RoundToEven(maxAge.Seconds() + 0.5))
		stsValue := fmt.Sprintf("max-age=%d; includeSubDomains", seconds)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if seconds >= 0 && seconds <= TwoYears {
				w.Header().Set("Strict-Transport-Security", stsValue)
			}

			next.ServeHTTP(w, r)
		})
	}
}
