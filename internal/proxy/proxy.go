package proxy

import (
	"log"
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

		// Inject custom sidecar version header
		req.Header.Set("X-Sidecar-Version", "1.0")

		// Add forwarding headers (these work)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Real-IP", req.RemoteAddr)

		// Remove sensitive headers from client request (works)
		req.Header.Del("Authorization")

		// Path rewriting (works)
		req.URL.Path = "/caleb" + req.URL.Path

		// Log all request headers
		log.Printf("=== Request Headers for %s %s ===", req.Method, req.URL.Path)
		for name, values := range req.Header {
			for _, value := range values {
				log.Printf("  %s: %s", name, value)
			}
		}
		log.Printf("=== End Headers ===")
	}

	// ModifyResponse lets you modify the response from upstream
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Add custom header (works)
		resp.Header.Set("X-Proxy", "opentelemetry-sidecar")

		// Remove upstream server identification (works)
		resp.Header.Del("Server")

		log.Printf("Response from upstream: %d %s", resp.StatusCode, resp.Status)
		return nil
	}

	return &ReverseProxy{
		proxy:    proxy,
		upstream: upstream,
	}, nil
}

// ServeHTTP forwards the request to the upstream server.
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Wrap the request body with CountingReader to track bytes without breaking streaming
	var countingReader *CountingReader
	if r.Body != nil {
		countingReader = NewCountingReader(r.Body)
		r.Body = countingReader
	}

	rp.proxy.ServeHTTP(w, r)

	// Log the byte count after the request is complete
	if countingReader != nil {
		log.Printf("Request body bytes read: %d", countingReader.BytesRead())
	}
}

// Upstream returns the upstream URL.
func (rp *ReverseProxy) Upstream() *url.URL {
	return rp.upstream
}
