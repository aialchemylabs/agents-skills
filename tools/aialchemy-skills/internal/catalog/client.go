package catalog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultMaxBytes bounds the catalog response size.
	DefaultMaxBytes = int64(4 << 20)
	// DefaultTimeout bounds the whole fetch, connection included.
	DefaultTimeout = 30 * time.Second
)

// FetchOptions tune catalog retrieval.
type FetchOptions struct {
	// AllowLoopbackHTTP permits plain HTTP for loopback hosts only. Callers
	// enable it only when the catalog URL was explicitly selected (flag or
	// environment), never for the built-in default.
	AllowLoopbackHTTP bool
	MaxBytes          int64
	Timeout           time.Duration
	// Client overrides the HTTP client (used by tests to trust test TLS certs).
	Client *http.Client
	APIKey string
}

// Fetch retrieves the catalog document with scheme, size, and time bounds.
// It makes no filesystem changes.
func Fetch(ctx context.Context, rawURL string, opts FetchOptions) ([]byte, error) {
	if opts.MaxBytes <= 0 {
		opts.MaxBytes = DefaultMaxBytes
	}
	if opts.Timeout <= 0 {
		opts.Timeout = DefaultTimeout
	}

	_, err := validateCatalogURL(rawURL, opts.AllowLoopbackHTTP)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid catalog request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	if opts.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+opts.APIKey)
	}
	req.Header.Set("User-Agent", "aialchemy-skills/0.5.0")

	client := opts.Client
	if client == nil {
		client = &http.Client{
			Timeout: opts.Timeout,
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: opts.Timeout,
				Proxy:                 http.ProxyFromEnvironment,
			},
		}
	}
	clientCopy := *client
	previousRedirectCheck := client.CheckRedirect
	clientCopy.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if req.URL.Scheme != via[0].URL.Scheme || req.URL.Host != via[0].URL.Host {
			return errors.New("cross-origin catalog redirects are not permitted")
		}
		if _, err := validateCatalogURL(req.URL.String(), opts.AllowLoopbackHTTP); err != nil {
			return err
		}
		if previousRedirectCheck != nil {
			return previousRedirectCheck(req, via)
		}
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	resp, err := clientCopy.Do(req)
	if err != nil {
		return nil, fmt.Errorf("catalog fetch failed: %s", Redact(err.Error()))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog fetch failed: HTTP %d from %s", resp.StatusCode, Redact(rawURL))
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, opts.MaxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("catalog read failed: %s", Redact(err.Error()))
	}
	if int64(len(data)) > opts.MaxBytes {
		return nil, fmt.Errorf("catalog response exceeds the %d-byte limit", opts.MaxBytes)
	}
	return data, nil
}

func validateCatalogURL(rawURL string, allowLoopbackHTTP bool) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid catalog URL: %v", err)
	}
	if u.Hostname() == "" {
		return nil, errors.New("catalog URL has no host")
	}
	if u.User != nil {
		return nil, errors.New("catalog URL must not contain embedded credentials")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("catalog URL must not contain a query or fragment")
	}
	switch u.Scheme {
	case "https":
	case "http":
		if !allowLoopbackHTTP || !isLoopbackHost(u.Hostname()) {
			return nil, fmt.Errorf("catalog URL %s: plain HTTP is allowed only for explicitly selected loopback test URLs", Redact(rawURL))
		}
	default:
		return nil, fmt.Errorf("catalog URL %s: unsupported scheme %q (use https)", Redact(rawURL), u.Scheme)
	}
	return u, nil
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
