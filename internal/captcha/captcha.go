package captcha

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	//go:embed captcha.html
	captchaHtml []byte
)

func New(ctx context.Context) (string, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	captchaUrl := fmt.Sprintf("http://127.0.0.1:%v", l.Addr().(*net.TCPAddr).Port)

	g, ctx := errgroup.WithContext(ctx)

	serverCtx, serverCancel := context.WithCancel(ctx)

	g.Go(func() error {
		return serve(serverCtx, l)
	})

	var ret string

	g.Go(func() error {
		defer serverCancel()
		var err error
		ret, err = execute(ctx, captchaUrl)
		return err
	})

	return ret, g.Wait()
}

func serve(ctx context.Context, l net.Listener) error {
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

func execute(ctx context.Context, captchaUrl string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	output, err := exec.CommandContext(ctx, "python3", path.Join(wd, "third_party", "hcaptcha-challenger", "src", "main.py"), "demo", captchaUrl, "--silence").Output()
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "CAPTCHA_RESPONSE: ") {
			return line[len("CAPTCHA_RESPONSE: "):], nil
		}
	}
	return "", errors.New("CAPTCHA_RESPONSE not found")
}
