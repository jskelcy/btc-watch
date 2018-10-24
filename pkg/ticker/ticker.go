package ticker

import (
	"fmt"
	"net/http"

	"github.com/jskelcy/btc-watch/pkg/aggregation"
	"github.com/jskelcy/btc-watch/pkg/fetch"
)

// Ticker implments functions compatible with http.HandleFunc.
// Functions are for interacting with aggregated BTC price data.
type Ticker interface {
	// Start starts consuming BTC price data.
	Start()
	// Price fetches current aggregated BTC price and writes it to http response.
	Price(http.ResponseWriter, *http.Request)
}

// Config is configuration for a new ticker.
type Config struct {
	Aggregator aggregation.Aggregator
	Fetcher    fetch.Fetcher
}

// NewTicker returns a Ticker from config.
func NewTicker(cfg Config) Ticker {
	return &ticker{
		aggregator: cfg.Aggregator,
		fetcher:    cfg.Fetcher,
	}
}

type ticker struct {
	aggregator aggregation.Aggregator
	fetcher    fetch.Fetcher
}

// Start starts consuming BTC price data.
func (t *ticker) Start() {
	priceChan := t.fetcher.Start()
	go func() {
		for {
			select {
			case price := <-priceChan:
				t.aggregator.Ingest(price)
			}
		}
	}()
}

// Price fetches current aggregated BTC price and writes it to http response.
func (t *ticker) Price(resp http.ResponseWriter, req *http.Request) {
	currAggPrice, err := t.aggregator.CurrentAggregatedPrice()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("1 minute moving average not yet available"))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf("%.2f", currAggPrice)))
}
