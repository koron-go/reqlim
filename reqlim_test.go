package reqlim

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func httpGet(t *testing.T, c *http.Client, url string) (int, string) {
	t.Helper()
	r, err := c.Get(url)
	if err != nil {
		t.Fatalf("failed to GET %s: %s", url, err)
	}
	var body string
	if r.Body != nil {
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %s", err)
		}
		body = string(b)
	}
	return r.StatusCode, body
}

func TestReqlim(t *testing.T) {
	s := httptest.NewServer(Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "OK")
		}), 2, ""))
	t.Cleanup(func() { s.Close() })

	c := s.Client()
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		sc, body := httpGet(t, c, s.URL)
		if sc != http.StatusOK || body != "OK" {
			t.Errorf("unexpected get#1: sc=%d body=%s", sc, body)
		}
	}()
	go func() {
		defer wg.Done()
		sc, body := httpGet(t, c, s.URL)
		if sc != http.StatusOK || body != "OK" {
			t.Errorf("unexpected get#2: sc=%d body=%s", sc, body)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		defer wg.Done()
		sc, body := httpGet(t, c, s.URL)
		if sc != http.StatusServiceUnavailable || body != defaultErrorBody {
			t.Errorf("unexpected get#3: sc=%d body=%s", sc, body)
		}
	}()
	wg.Wait()
}

func TestReqlim_body(t *testing.T) {
	s := httptest.NewServer(Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "OK")
		}), 2, "BUSY"))
	t.Cleanup(func() { s.Close() })

	c := s.Client()
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		sc, body := httpGet(t, c, s.URL)
		if sc != http.StatusOK || body != "OK" {
			t.Errorf("unexpected get#1: sc=%d body=%s", sc, body)
		}
	}()
	go func() {
		defer wg.Done()
		sc, body := httpGet(t, c, s.URL)
		if sc != http.StatusOK || body != "OK" {
			t.Errorf("unexpected get#2: sc=%d body=%s", sc, body)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		defer wg.Done()
		sc, body := httpGet(t, c, s.URL)
		if sc != http.StatusServiceUnavailable || body != "BUSY" {
			t.Errorf("unexpected get#3: sc=%d body=%s", sc, body)
		}
	}()
	wg.Wait()
}
