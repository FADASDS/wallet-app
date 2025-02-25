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

func TestGetLastNTransactions_ok(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.TransactionHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/transactions?count=3", nil)
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
	var respData response.DataOkResponse[[]dto.TransactionDTO]
	err = json.Unmarshal(body, &respData)

	if err != nil {
		t.Fatal(err)
	}

	if respData.Data == nil {
		t.Fatal("response data is nil")
	}

	if len(*respData.Data) != 3 {
		t.Error("invalid response data")
	}
}

func TestGetLastNTransactions_invalidCount(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.TransactionHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/transactions?count=-3", nil)
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

	if errResp.Error != "Invalid 'count' parameter." {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "Invalid 'count' parameter.")
	}
}

func TestGetLastNTransactions_noParams(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.TransactionHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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

	if errResp.Error != "'count' parameter is missing." {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "'count' parameter is missing.")
	}
}

func TestGetLastNTransactions_wrongMethod(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.TransactionHandler{Store: store}

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", nil)
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
