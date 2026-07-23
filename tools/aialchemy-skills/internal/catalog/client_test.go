package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchLoopbackHTTPOnlyInTestMode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"plugins": []}`))
	}))
	defer ts.Close()

	if _, err := Fetch(context.Background(), ts.URL, FetchOptions{}); err == nil {
		t.Fatal("loopback HTTP accepted without explicit test mode")
	}
	data, err := Fetch(context.Background(), ts.URL, FetchOptions{AllowLoopbackHTTP: true})
	if err != nil {
		t.Fatalf("explicit loopback HTTP fetch failed: %v", err)
	}
	if !strings.Contains(string(data), "plugins") {
		t.Fatalf("unexpected body: %s", data)
	}
}

func TestFetchRejectsNonLoopbackHTTP(t *testing.T) {
	// Scheme/host validation happens before any dial, so no request is made.
	_, err := Fetch(context.Background(), "http://example.com/catalog.json", FetchOptions{AllowLoopbackHTTP: true})
	if err == nil || !strings.Contains(err.Error(), "loopback") {
		t.Fatalf("err = %v, want loopback restriction", err)
	}
}

func TestFetchRejectsBadSchemes(t *testing.T) {
	for _, u := range []string{
		"ftp://example.com/x",
		"file:///etc/passwd",
		"https://user:secret@example.com/catalog.json",
		"https://example.com/catalog.json?token=secret",
		"https://example.com/catalog.json#fragment",
	} {
		if _, err := Fetch(context.Background(), u, FetchOptions{}); err == nil {
			t.Errorf("scheme accepted: %s", u)
		}
	}
}

func TestFetchRejectsInsecureRedirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://example.com/catalog.json", http.StatusFound)
	}))
	defer ts.Close()
	_, err := Fetch(context.Background(), ts.URL, FetchOptions{AllowLoopbackHTTP: true})
	if err == nil || !strings.Contains(err.Error(), "loopback") {
		t.Fatalf("redirect error = %v, want loopback HTTP refusal", err)
	}
}

func TestFetchHTTPS(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"plugins": []}`))
	}))
	defer ts.Close()
	data, err := Fetch(context.Background(), ts.URL, FetchOptions{Client: ts.Client()})
	if err != nil {
		t.Fatalf("https fetch failed: %v", err)
	}
	if !strings.Contains(string(data), "plugins") {
		t.Fatalf("unexpected body: %s", data)
	}
}

func TestFetchOversizedResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strings.Repeat("x", 200)))
	}))
	defer ts.Close()
	_, err := Fetch(context.Background(), ts.URL, FetchOptions{AllowLoopbackHTTP: true, MaxBytes: 100})
	if err == nil || !strings.Contains(err.Error(), "limit") {
		t.Fatalf("err = %v, want size-limit error", err)
	}
}

func TestFetchTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer ts.Close()
	start := time.Now()
	_, err := Fetch(context.Background(), ts.URL, FetchOptions{AllowLoopbackHTTP: true, Timeout: 150 * time.Millisecond})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("timeout took %v, want < 1s", elapsed)
	}
}

func TestFetchNon200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "gone", http.StatusNotFound)
	}))
	defer ts.Close()
	_, err := Fetch(context.Background(), ts.URL, FetchOptions{AllowLoopbackHTTP: true})
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("err = %v, want HTTP 404", err)
	}
}
