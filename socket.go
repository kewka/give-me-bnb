package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

func NewSocketTransaction(
	ctx context.Context,
	socketUrl string,
	captcha string,
	account string,
	proxy string,
) (string, error) {
	dialer := websocket.Dialer{}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		dialer.Proxy = http.ProxyURL(proxyUrl)
	}
	conn, _, err := dialer.Dial(socketUrl, nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return requestSocketBnb(conn, captcha, account)
	})

	var ret string
	g.Go(func() error {
		var err error
		ret, err = waitSocketTransaction(ctx, conn, account)
		return err
	})
	return ret, g.Wait()
}

func requestSocketBnb(conn *websocket.Conn, captcha string, account string) error {
	return conn.WriteJSON(map[string]interface{}{
		"url":     account,
		"symbol":  "BNB",
		"tier":    0,
		"captcha": captcha,
	})
}

type SocketMessage struct {
	Error    *string `json:"error"`
	Requests []struct {
		Account string `json:"account"`
		Tx      struct {
			Hash string `json:"hash"`
		} `json:"tx"`
	} `json:"requests"`
}

func waitSocketTransaction(ctx context.Context, conn *websocket.Conn, account string) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			msg := &SocketMessage{}
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
