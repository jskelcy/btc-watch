package main

import (
	"log"

	"github.com/jskelcy/btc-cli/pkg/aggregation"
	"github.com/jskelcy/btc-cli/pkg/fetch"
	"github.com/jskelcy/btc-cli/pkg/server"
	"github.com/jskelcy/btc-cli/pkg/ticker"
)

func main() {
	aggregator := aggregation.NewAggregator(aggregation.Config{
		AggWindow:        60,
		CollectionWindow: 1,
	})
	fetcher := fetch.NewFetcher(fetch.Config{
		EndpointBase:   "https://api.binance.com",
		PricesEndpoint: "api/v3/ticker/price",
		Symbol:         "BTCUSDT",
	})
	priceTicker := ticker.NewTicker(ticker.Config{
		Aggregator: aggregator,
		Fetcher:    fetcher,
	})
	priceTicker.Start()
	server := server.NewServer(priceTicker)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
