package main

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed captcha.html
	captchaHtml []byte
)

func NewCaptcha(ctx context.Context) (string, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	captchaUrl := fmt.Sprintf("http://localhost:%v", l.Addr().(*net.TCPAddr).Port)

	g, ctx := errgroup.WithContext(ctx)

	serverCtx, serverCancel := context.WithCancel(ctx)
	defer serverCancel()

	g.Go(func() error {
		return serveCaptcha(serverCtx, l)
	})

	var ret string

	g.Go(func() error {
		defer serverCancel()
		var err error
		ret, err = executeCaptcha(ctx, captchaUrl)
		return err
	})

	return ret, g.Wait()
}

func serveCaptcha(ctx context.Context, l net.Listener) error {
	g, ctx := errgroup.WithContext(ctx)
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write(captchaHtml)
		}),
	}

	g.Go(func() error {
		err := srv.Serve(l)
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	})

	g.Go(func() error {
		// Handle shutdown
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	return g.Wait()
}

func executeCaptcha(ctx context.Context, captchaUrl string) (string, error) {
	var ret string
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate(captchaUrl),
		chromedp.Poll("grecaptcha", nil),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var promise *runtime.RemoteObject
			err := chromedp.Evaluate(
				"new Promise(resolve => grecaptcha.execute().then(resolve))",
				&promise,
			).Do(ctx)
			if err != nil {
				return err
			}

			res, _, err := runtime.AwaitPromise(promise.ObjectID).Do(ctx)
			if err != nil {
				return err
			}

			ret = string(res.Value)
			// remove quoutes
			ret = ret[1 : len(ret)-1]
			return nil
		}),
	)
	return ret, err
}
