package proxy

import (
	"log"
	"net/http"
)

const (
	// DefaultListenAddr is the default address the proxy listens on.
	DefaultListenAddr = ":8080"
	// DefaultUpstreamURL is the default upstream server URL.
	DefaultUpstreamURL = "http://localhost:3000"
)

// Handler wraps the reverse proxy and provides HTTP handler functionality.
type Handler struct {
	proxy      *ReverseProxy
	listenAddr string
}

// NewHandler creates a new proxy handler with the given configuration.
func NewHandler(listenAddr, upstreamURL string) (*Handler, error) {
	if listenAddr == "" {
		listenAddr = DefaultListenAddr
	}
	if upstreamURL == "" {
		upstreamURL = DefaultUpstreamURL
	}

	proxy, err := NewReverseProxy(upstreamURL)
	if err != nil {
		return nil, err
	}

	return &Handler{
		proxy:      proxy,
		listenAddr: listenAddr,
	}, nil
}

// ServeHTTP implements http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Proxying request: %s %s -> %s", r.Method, r.URL.Path, h.proxy.Upstream())
	h.proxy.ServeHTTP(w, r)
}

// ListenAndServe starts the HTTP server on the configured listen address.
func (h *Handler) ListenAndServe() error {
	log.Printf("Starting proxy server on %s, forwarding to %s", h.listenAddr, h.proxy.Upstream())
	return http.ListenAndServe(h.listenAddr, h)
}

// ListenAddr returns the configured listen address.
func (h *Handler) ListenAddr() string {
	return h.listenAddr
}
