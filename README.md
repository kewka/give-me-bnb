# give-me-bnb

Automation for https://testnet.binance.org/faucet-smart with some hacks 😈

## Usage

```sh
$ give-me-bnb -help
Usage of give-me-bnb:
  -proxy string
        proxy url
  -rpc-url string
        bsc testnet rpc url (default "https://data-seed-prebsc-1-s1.binance.org:8545")
  -socket-url string
        bsc faucet socket url (default "wss://testnet.binance.org/faucet-smart/api")
  -to string
        your address (required)
```

## Installation

### Go binary

```sh
$ go install github.com/kewka/give-me-bnb@latest
```

### Docker image

The docker image has built-in [Tor](https://www.torproject.org/) proxy server ([see example below](#docker-tor-proxy)).

```sh
$ docker pull kewka/give-me-bnb
```

## Examples

### Binary

```sh
$ give-me-bnb -to <address>
```

### Docker

```sh
$ docker run --rm kewka/give-me-bnb give-me-bnb -to <address>
```

### Docker (Tor proxy)

```sh
$ docker run --rm kewka/give-me-bnb give-me-bnb -proxy socks5://127.0.0.1:9050 -to <address>
```

### Docker (Tor proxy + infinite loop)

```sh
$ docker run --rm kewka/give-me-bnb sh -c "while :; do give-me-bnb -proxy socks5://127.0.0.1:9050 -to <address>; killall -HUP tor; done"
```
