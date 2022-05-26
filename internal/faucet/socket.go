package faucet

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

func NewTransaction(
	ctx context.Context,
	captcha string,
	account string,
	proxy string,
	symbol string,
) (string, error) {
	dialer := websocket.Dialer{}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		dialer.Proxy = http.ProxyURL(proxyUrl)
	}
	conn, _, err := dialer.Dial("wss://testnet.binance.org/faucet-smart/api", nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return claim(conn, captcha, account, symbol)
	})

	var ret string
	g.Go(func() error {
		var err error
		ret, err = waitTransaction(ctx, conn, account)
		return err
	})
	return ret, g.Wait()
}

func claim(conn *websocket.Conn, captcha string, account string, symbol string) error {
	return conn.WriteJSON(map[string]interface{}{
		"url":     account,
		"symbol":  symbol,
		"tier":    0,
		"captcha": captcha,
	})
}

type Message struct {
	Error    *string `json:"error"`
	Requests []struct {
		Account string `json:"account"`
		Tx      struct {
			Hash string `json:"hash"`
		} `json:"tx"`
	} `json:"requests"`
}

func waitTransaction(ctx context.Context, conn *websocket.Conn, account string) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			msg := Message{}
			if err := conn.ReadJSON(&msg); err != nil {
				return "", err
			}

			if msg.Error != nil {
				return "", fmt.Errorf("faucet error: %v", *msg.Error)
			}

			for _, r := range msg.Requests {
				if strings.EqualFold(r.Account, account) {
					return r.Tx.Hash, nil
				}
			}
		}
	}
}
