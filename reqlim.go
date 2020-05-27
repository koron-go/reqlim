package reqlim

import (
	"io"
	"net/http"

	"golang.org/x/sync/semaphore"
)

// Reqlim limits number of request in parallel processing.
type Reqlim struct {
	handler http.Handler
	sem     *semaphore.Weighted

	body    string
}

// Handler returns a http.Handler that runs h with the given concurrency limit.
func Handler(h http.Handler, limit int, msg string) http.Handler {
	return &Reqlim{
		sem:     semaphore.NewWeighted(int64(limit)),
		handler: h,
		body:    msg,
	}
}

// ServeHTTP implements http.Handler
func (rl *Reqlim) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !rl.sem.TryAcquire(1) {
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, rl.errorBody())
		return
	}
	rl.handler.ServeHTTP(w, r)
	rl.sem.Release(1)
}

func (rl *Reqlim) errorBody() string {
	if rl.body != "" {
		return rl.body
	}
	return defaultErrorBody
}

const defaultErrorBody = "<html><head><title>Server busy</title></head><body><h1>Server busy</h1></body></html>"
