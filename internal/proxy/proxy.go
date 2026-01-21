package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy wraps httputil.ReverseProxy with upstream configuration.
type ReverseProxy struct {
	proxy    *httputil.ReverseProxy
	upstream *url.URL
}

// NewReverseProxy creates a new reverse proxy that forwards requests
// to the specified upstream URL.
func NewReverseProxy(upstreamURL string) (*ReverseProxy, error) {
	upstream, err := url.Parse(upstreamURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(upstream)

	// Customize the Director to ensure Host header is set correctly
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = upstream.Host
	}

	return &ReverseProxy{
		proxy:    proxy,
		upstream: upstream,
	}, nil
}

// ServeHTTP forwards the request to the upstream server.
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rp.proxy.ServeHTTP(w, r)
}

// Upstream returns the upstream URL.
func (rp *ReverseProxy) Upstream() *url.URL {
	return rp.upstream
}
