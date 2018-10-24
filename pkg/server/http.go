package server

import (
	"fmt"
	"net/http"

	"github.com/jskelcy/btc-watch/pkg/ticker"
)

const (
	portFmt = ":%s"
)

// Config contians configuration for an HTTP server.
type Config struct {
	PriceTicker ticker.Ticker
	Port        string
}

// NewServer returns a new HTTP server from config.
func NewServer(cfg Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/price", handleMethod(http.MethodGet, cfg.PriceTicker.Price))

	return &http.Server{
		Addr:    fmt.Sprintf(portFmt, cfg.Port),
		Handler: mux,
	}
}

func handleMethod(method string, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "text/plain")
			resp.Write([]byte("HTTP method not supported"))
			return
		}

		handler(resp, req)
	}
}
