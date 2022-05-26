package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kewka/give-me-bnb/internal/blockchain"
	"github.com/kewka/give-me-bnb/internal/captcha"
	"github.com/kewka/give-me-bnb/internal/faucet"
)

var (
	proxy    string
	rpcUrl   string
	account  *blockchain.Account
	currency faucet.Currency
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func parseArgs() error {
	keyArg := flag.String("key", "", "your private key (required)")
	flag.StringVar(&proxy, "proxy", "", "proxy url")
	flag.StringVar(&rpcUrl, "rpc", "https://data-seed-prebsc-1-s1.binance.org:8545", "bsc testnet rpc url")
	currencyArg := flag.String("currency", faucet.BNB.Symbol(), "faucet currency")
	flag.Parse()

	if *keyArg == "" {
		flag.Usage()
		return errors.New("`key` cannot be empty")
	}

	var err error

	account, err = blockchain.ReadAccount(*keyArg)
	if err != nil {
		return fmt.Errorf("blockchain.ReadAccount: %w", err)
	}

	currency, err = faucet.NewCurrency(*currencyArg)
	if err != nil {
		return fmt.Errorf("faucet.NewCurrency: %w", err)
	}

	return nil
}

func run() error {
	if err := parseArgs(); err != nil {
		return err
	}

	fmt.Println("account:", account.PublicHex())
	fmt.Println("currency:", currency.Symbol())

	baseCtx := context.Background()

	client, err := blockchain.NewClient(rpcUrl)
	if err != nil {
		return fmt.Errorf("blockchain.NewClient: %w", err)
	}
	defer client.Close()

	captchaRes, err := captcha.New(baseCtx)
	if err != nil {
		return fmt.Errorf("captcha.New: %w", err)
	}

	tempAccount, err := blockchain.GenerateAccount()
	if err != nil {
		return fmt.Errorf("blockchain.GenerateAccount: %w", err)
	}

	fmt.Printf("temp: %v (%v)\n", tempAccount.PublicHex(), tempAccount.PrivateHex())

	ctx, cancel := context.WithTimeout(baseCtx, 2*time.Minute)
	defer cancel()
	faucetTx, err := faucet.NewTransaction(
		ctx,
		captchaRes,
		tempAccount.PublicHex(),
		proxy,
		currency.Symbol(),
	)
	if err != nil {
		return fmt.Errorf("faucet.NewTransaction: %w", err)
	}

	fmt.Println("faucet->temp:", faucetTx)

	ctx, cancel = context.WithTimeout(baseCtx, time.Minute)
	defer cancel()
	if err := client.Wait(ctx, faucetTx); err != nil {
		return fmt.Errorf("client.Wait: %w", err)
	}

	ctx, cancel = context.WithTimeout(baseCtx, time.Minute)
	defer cancel()

	var tx string

	currencyAddr := currency.Address()
	if currencyAddr != "" {
		tx, err = client.WithdrawToken(ctx, tempAccount, account, currencyAddr)
		if err != nil {
			return fmt.Errorf("client.WithdrawToken: %w", err)
		}
	} else {
		tx, err = client.WithdrawNative(ctx, tempAccount, account)
		if err != nil {
			return fmt.Errorf("client.WithdrawNative: %w", err)
		}
	}
	fmt.Println("temp->account:", tx)
	return nil
}
