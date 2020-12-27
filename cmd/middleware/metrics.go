//go:generate mockgen -destination=../../tests/mocks/mock_middleware.go -package=mocks -source middleware.go

package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//NewInboundObserver uses a reporter to create a new inbound request middleware the provides visibility for inbound calls.
//reports the response time and http status code for all app api paths.
func NewInboundObserver(reporter reporter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			responseRecorder := newResponseRecorder(w)
			path, _ := mux.CurrentRoute(r).GetPathTemplate()
			defer func() {
				reporter.ObserverHistogram("inbound_request_histogram", float64(time.Since(start).Milliseconds()), r.Method, path, strconv.Itoa(responseRecorder.statusCode))
				reporter.ObserverSummary("inbound_request_summary", float64(time.Since(start).Milliseconds()), r.Method, path, strconv.Itoa(responseRecorder.statusCode))
			}()

			next.ServeHTTP(responseRecorder, r)
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{w, http.StatusOK}
}

func (lrw *responseRecorder) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
