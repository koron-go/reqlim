package reqlim

import (
	"net/http"

	"golang.org/x/sync/semaphore"
)

// Reqlim limits number of request in parallel processing.
type Reqlim struct {
	s *semaphore.Weighted
	h http.Handler
}

// New creates a request limitter with weight `n`.
func New(n int64, h http.Handler) *Reqlim {
	return &Reqlim{
		s: semaphore.NewWeighted(n),
	}
}

// ServeHTTP implements http.Handler
func (rl *Reqlim) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !rl.s.TryAcquire(1) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("server busy"))
		return
	}
	rl.h.ServeHTTP(w, r)
	rl.s.Release(1)
}
