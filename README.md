# give-me-bnb

Automation for https://testnet.binance.org/faucet-smart with some hacks 😈

## Usage

```sh
$ give-me-bnb -help
Usage of give-me-bnb:
  -currency string
        faucet currency (default "BNB")
  -key string
        your private key (required)
  -proxy string
        proxy url
  -rpc string
        bsc testnet rpc url (default "https://data-seed-prebsc-1-s1.binance.org:8545")
```

## Installation

### Docker image

The docker image has built-in [Tor](https://www.torproject.org/) proxy server ([see example below](#docker-tor-proxy)).

```sh
$ docker pull kewka/give-me-bnb
```

## Examples

### Binary

```sh
$ give-me-bnb -key <key>
```

### Docker

```sh
$ docker run --rm kewka/give-me-bnb give-me-bnb -key <key>
```

### Docker (Tor proxy)

```sh
$ docker run --rm kewka/give-me-bnb give-me-bnb -proxy socks5://127.0.0.1:9050 -key <key>
```

### Docker (Tor proxy + infinite loop)

```sh
$ docker run --rm kewka/give-me-bnb sh -c "while :; do give-me-bnb -proxy socks5://127.0.0.1:9050 -key <key>; killall -HUP tor; done"
```

## Credits

- https://github.com/QIN2DIM/hcaptcha-challenger
