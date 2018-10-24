# btc-watch

btc-watch fetches BTC prices in USD every second and keeps a 1 minute moving average.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Go 1.11.1

### Building and running

To build this project use the `build` make target to build the binary in `/out` directory. When using `build`, dependecies are also installed.

```
make build
```

Finally you can install dependencies, build a binary, and run the binary with the `run` make target.
The default port is `8080` but that can be upadted using `port` flag.

```
make run port=8081
```

## Running the tests

Tests can be run using the `test` make target.

```
make test
```

## API

btc-watch offers only one HTTP GET endpoint `/price`. This endpoint returns the 1 minute moving average of the price of BTC in USD.

For example:
```
curl localhost:8080/price
6597.57
```