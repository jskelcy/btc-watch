package fetch

import (
	"net/http"
	"testing"
	"time"
)

func TestFetchWithNoErrors(t *testing.T) {
	f := &fetcher{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		symbol:         "BTCUSDT",
		endpointBase:   "https://api.binance.com",
		pricesEndpoint: "api/v3/ticker/price",
	}

	_, err := f.getPrice()
	if err != nil {
		t.Error(err)
	}
}
