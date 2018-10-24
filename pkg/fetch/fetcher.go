package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	endpointFormat = "%s/%s"
	symbolQueryArg = "symbol"
)

// Fetcher fetches currency price data at some interval and
// provides streaming updates of a channel.
type Fetcher interface {
	// Start will start fetching currency price updates at some interval.
	Start() <-chan float64
}

// Config is the configuration for a new Fetcher.
type Config struct {
	// PricesEndpoint is the endpoint where prices are served from the provider.
	PricesEndpoint string
	// EndpointBase is the root/provider of prices API.
	EndpointBase string
	// Symbol is the symbol for the currency to be fetched
	Symbol string
}

// NewFetcher returns a Fetcher from config.
func NewFetcher(cfg Config) Fetcher {
	return &fetcher{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		endpointBase:   cfg.EndpointBase,
		priceChan:      make(chan float64),
		pricesEndpoint: cfg.PricesEndpoint,
		symbol:         cfg.Symbol,
	}
}

type fetcher struct {
	client         *http.Client
	symbol         string
	priceChan      chan float64
	endpointBase   string
	pricesEndpoint string
}

// Start will start fetching currency price updates at some interval.
func (f *fetcher) Start() <-chan float64 {
	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				price, err := f.getPrice()
				if err != nil {
					log.Panicf("error getting price: %v", err.Error())
				}
				f.priceChan <- price
			}
		}
	}()

	return f.priceChan
}

func (f *fetcher) getPrice() (float64, error) {
	req, err := f.buildQuery()
	if err != nil {
		return 0, err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return 0, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	pr := &priceRespBody{}
	err = json.Unmarshal(data, pr)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(pr.Price, 64)
}

func (f *fetcher) buildQuery() (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(endpointFormat, f.endpointBase, f.pricesEndpoint),
		nil,
	)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add(symbolQueryArg, f.symbol)
	req.URL.RawQuery = q.Encode()

	return req, nil
}
