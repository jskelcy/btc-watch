package server

import (
	"net/http"

	"github.com/jskelcy/btc-cli/pkg/ticker"
)

func NewServer(priceTicker ticker.Ticker) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/price", handleMethod(http.MethodGet, priceTicker.Price))

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func handleMethod(method string, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte("HTTP method not supported"))
			return
		}

		handler(resp, req)
	}
}
