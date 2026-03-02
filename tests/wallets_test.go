package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gersastas/wallet-service/internal/transport/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	srv := server.New(":1010")

	testServer := httptest.NewServer(srv.Handler())

	t.Cleanup(func() {
		testServer.Close()
	})

	return testServer
}

func makeRequest(t *testing.T, ts *httptest.Server, method, path string, body interface{}) *http.Response {
	t.Helper()

	var reqBody bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&reqBody).Encode(body)
		require.NoError(t, err, "failed to encode request body")
	}

	req, err := http.NewRequest(method, ts.URL+path, &reqBody)
	require.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", "application/json")

	resp, err := ts.Client().Do(req)
	require.NoError(t, err, "failed to send request")

	return resp
}

func decodeResponse(t *testing.T, resp *http.Response, target interface{}) {
	t.Helper()

	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err, "failed to close response body")
	}()

	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err, "failed to decode response")
}

func TestCreateWallet_Success(t *testing.T) {
	ts := setupTestServer(t)
	reqBody := server.WalletRequest{
		UserID:   "550e8400-e29b-41d4-a716-446655440000",
		Name:     "Test Wallet",
		Currency: "RUB",
	}

	resp := makeRequest(t, ts, http.MethodPost, "/wallets", reqBody) // ← убрал слэш

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result server.WalletResponse
	decodeResponse(t, resp, &result)

	assert.NotEmpty(t, result.ID)
	assert.Equal(t, reqBody.UserID, result.UserID)
	assert.Equal(t, reqBody.Name, result.Name)
	assert.Equal(t, int64(0), result.Balance)
}

func TestCreateWallet_ValidationError(t *testing.T) {
	ts := setupTestServer(t)

	testCases := []struct {
		name    string
		body    server.WalletRequest
		wantErr string
	}{
		{
			name:    "empty user_id",
			body:    server.WalletRequest{Name: "Test", Currency: "RUB"},
			wantErr: "user_id is required",
		},
		{
			name:    "empty name",
			body:    server.WalletRequest{UserID: "550e8400-e29b-41d4-a716-446655440000", Currency: "RUB"},
			wantErr: "name is required",
		},
		{
			name:    "empty currency",
			body:    server.WalletRequest{UserID: "550e8400-e29b-41d4-a716-446655440000", Name: "Test"},
			wantErr: "currency is required",
		},
		{
			name:    "invalid user_id format",
			body:    server.WalletRequest{UserID: "invalid-uuid", Name: "Test", Currency: "RUB"},
			wantErr: "user_id is invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := makeRequest(t, ts, http.MethodPost, "/wallets", tc.body) // ← убрал слэш

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var errResp server.ErrorResponse
			decodeResponse(t, resp, &errResp)
			assert.Contains(t, errResp.Error, tc.wantErr)
		})
	}
}

func TestGetWallet_Success(t *testing.T) {
	ts := setupTestServer(t)

	createReq := server.WalletRequest{
		UserID:   "550e8400-e29b-41d4-a716-446655440000",
		Name:     "My Wallet",
		Currency: "USD",
	}
	createResp := makeRequest(t, ts, http.MethodPost, "/wallets", createReq) // ← убрал слэш

	var createdWallet server.WalletResponse
	decodeResponse(t, createResp, &createdWallet)
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	getResp := makeRequest(t, ts, http.MethodGet, "/wallets/"+createdWallet.ID, nil)

	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var result server.WalletResponse
	decodeResponse(t, getResp, &result)

	assert.Equal(t, createdWallet.ID, result.ID)
	assert.Equal(t, createdWallet.UserID, result.UserID)
	assert.Equal(t, createdWallet.Name, result.Name)
	assert.Equal(t, createdWallet.Balance, result.Balance)
	assert.Equal(t, createdWallet.Currency, result.Currency)
}

func TestGetWallet_NotFound(t *testing.T) {
	ts := setupTestServer(t)

	resp := makeRequest(t, ts, http.MethodGet, "/wallets/00000000-0000-0000-0000-000000000000", nil)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
