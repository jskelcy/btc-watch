package main

import (
	"flag"
	"log"

	"github.com/jskelcy/btc-watch/pkg/aggregation"
	"github.com/jskelcy/btc-watch/pkg/fetch"
	"github.com/jskelcy/btc-watch/pkg/server"
	"github.com/jskelcy/btc-watch/pkg/ticker"
)

const (
	defaultPort = "8080"
)

func main() {
	port := flag.String("port", defaultPort, "port to listen on")
	flag.Parse()

	if *port == "" {
		*port = defaultPort
	}

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
	server := server.NewServer(server.Config{
		PriceTicker: priceTicker,
		Port:        *port,
	})

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
