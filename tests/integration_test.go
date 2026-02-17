package tests

import (
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gersastas/wallet-service/internal/config"
	httpserver "github.com/gersastas/wallet-service/internal/transport/http/server"
	"github.com/sirupsen/logrus"
)

func TestServer_Integration(t *testing.T) {
	originalAddr := os.Getenv("HTTP_BIND_ADDR")
	defer func() {
		if err := os.Setenv("HTTP_BIND_ADDR", originalAddr); err != nil {
			t.Logf("failed to restore original HTTP_BIND_ADDR: %v", err)
		}
	}()

	port, err := getFreePort()
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	testAddr := "localhost:" + port
	if err := os.Setenv("HTTP_BIND_ADDR", testAddr); err != nil {
		t.Fatalf("failed to set HTTP_BIND_ADDR: %v", err)
	}

	cfg := config.New()
	if cfg.GetHTTPBindAddr() != testAddr {
		t.Fatalf("expected bind addr %q, got %q", testAddr, cfg.GetHTTPBindAddr())
	}

	server := httpserver.New(cfg.GetHTTPBindAddr())

	ready := make(chan struct{})

	go func() {
		close(ready)
		if err := server.Run(); err != nil {
			logrus.Errorf("server run failed: %v", err)
		}
	}()

	select {
	case <-ready:
		// Небольшая задержка для полной инициализации
		time.Sleep(100 * time.Millisecond)
	case <-time.After(2 * time.Second):
		t.Fatal("server did not start in time")
	}

	client := &http.Client{Timeout: 2 * time.Second}
	var resp *http.Response
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		resp, err = client.Get("http://" + testAddr + "/")
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("failed to connect to server after %d retries: %v", maxRetries, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close response body: %v", err)
		}
	}()

	t.Logf("server responded with status %s", resp.Status)
}

func getFreePort() (string, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := l.Close(); err != nil {
			_ = err
		}
	}()
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return "", err
	}
	return port, nil
}
