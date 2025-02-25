package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-app/internal/dto"
	"wallet-app/internal/http-server/handlers"
	"wallet-app/internal/lib/api/response"
	"wallet-app/internal/storage"
)

func TestGetWalletBalance_ok(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.WalletHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/wallet/wallet_1/balance", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading body")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	var respData response.DataOkResponse[dto.WalletBalanceDTO]
	err = json.Unmarshal(body, &respData)

	if err != nil {
		t.Fatal(err)
	}

	if respData.Data == nil {
		t.Fatal("response data is nil")
	}

	if respData.Data.WalletId != "wallet_1" {
		t.Error("invalid response balance")
	}

	if respData.Data.Balance != 100.0 {
		t.Error("invalid response balance")
	}
}

func TestGetWalletBalance_wrongWalletAddress(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.WalletHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/wallet/f/balance", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	var errResp response.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatal("JSON unmarshal error:", err)
	}

	if errResp.Error != "GetWalletBalance error!" {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "GetWalletBalance error!")
	}
}

func TestGetWalletBalance_noParams(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.WalletHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/wallet//balance", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	var errResp response.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatal("JSON unmarshal error:", err)
	}

	if errResp.Error != "Not enough parameters (address)!" {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "Not enough parameters (address)!")
	}
}

func TestGetWalletBalance_wrongMethod(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.WalletHandler{Store: store}

	req := httptest.NewRequest(http.MethodPost, "/api/wallet/wallet_1/balance", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	var errResp response.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatal("JSON unmarshal error:", err)
	}

	if errResp.Error != "Unsupported method!" {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "Unsupported method!")
	}
}
