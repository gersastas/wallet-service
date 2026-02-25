package tests

import (
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gersastas/wallet-service/internal/config"
	httpserver "github.com/gersastas/wallet-service/internal/transport/http/server"
	"github.com/stretchr/testify/require"
)

func TestServer_Integration(t *testing.T) {
	originalAddr := os.Getenv("HTTP_BIND_ADDR")
	defer func() {
		err := os.Setenv("HTTP_BIND_ADDR", originalAddr)
		require.NoError(t, err)
	}()

	port, err := getFreePort()
	require.NoError(t, err)

	testAddr := "localhost:" + port
	err = os.Setenv("HTTP_BIND_ADDR", testAddr)
	require.NoError(t, err)

	cfg := config.New()
	require.Equal(t, testAddr, cfg.GetHTTPBindAddr())

	server := httpserver.New(cfg.GetHTTPBindAddr())

	ready := make(chan struct{})

	go func() {
		close(ready)
		_ = server.Run()
	}()

	select {
	case <-ready:
		// Небольшая задержка для полной инициализации
		time.Sleep(100 * time.Millisecond)
	case <-time.After(2 * time.Second):
		t.Fatal("server did not start in time")
	}

	client := &http.Client{Timeout: 2 * time.Second}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := client.Get("http://" + testAddr + "/time")
			require.NoError(t, err)
			defer func() {
				err := resp.Body.Close()
				require.NoError(t, err)
			}()

			require.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}

	wg.Wait()

	t.Log("all 100 requests completed successfully")
}

func getFreePort() (string, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = l.Close()
	}()
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return "", err
	}
	return port, nil
}
