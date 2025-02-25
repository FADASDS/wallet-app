package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-app/internal/http-server/handlers"
	"wallet-app/internal/lib/api"
	"wallet-app/internal/lib/api/response"
	"wallet-app/internal/storage"
)

func TestSend_ok(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.SendHandler{Store: store}

	requestBody, err := json.Marshal(map[string]interface{}{
		"from":   "wallet_1",
		"to":     "wallet_2",
		"amount": 100.0,
	})

	if err != nil {
		t.Fatal("Failed to marshal request body:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/send", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()
	handler.Send(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading body")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	var respData response.EmptyOkResponse
	err = json.Unmarshal(body, &respData)

	if err != nil {
		t.Fatal(err)
	}

	if respData.Status != http.StatusOK {
		t.Error("invalid response status")
	}

}

func TestSend_wrongWalletAddress(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.SendHandler{Store: store}

	requestBody, err := json.Marshal(map[string]interface{}{
		"from":   "",
		"to":     "wallet_2",
		"amount": 100.0,
	})

	if err != nil {
		t.Fatal("Failed to marshal request body:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/send", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()
	handler.Send(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	var errResp api.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatal("JSON unmarshal error:", err)
	}

	if errResp.Error != "Internal server error." {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "GetWalletBalance error!")
	}
}

func TestSend_notEnoughBalance(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.SendHandler{Store: store}

	requestBody, err := json.Marshal(map[string]interface{}{
		"from":   "wallet_1",
		"to":     "wallet_2",
		"amount": 200.0,
	})

	if err != nil {
		t.Fatal("Failed to marshal request body:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/send", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()
	handler.Send(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	var errResp api.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatal("JSON unmarshal error:", err)
	}

	if errResp.Error != "Internal server error." {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "Internal server error.")
	}
}

func TestSend_invalidAmount(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.SendHandler{Store: store}

	requestBody, err := json.Marshal(map[string]interface{}{
		"from":   "wallet_1",
		"to":     "wallet_2",
		"amount": -200.0,
	})

	if err != nil {
		t.Fatal("Failed to marshal request body:", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/send", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()
	handler.Send(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	var errResp api.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatal("JSON unmarshal error:", err)
	}

	if errResp.Error != "Invalid amount" {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "Invalid amount")
	}
}

func TestSend_wrongMethod(t *testing.T) {
	store := &storage.MockStore{}
	handler := handlers.SendHandler{Store: store}

	requestBody, err := json.Marshal(map[string]interface{}{
		"from":   "wallet_1",
		"to":     "wallet_2",
		"amount": 100.0,
	})

	if err != nil {
		t.Fatal("Failed to marshal request body:", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/send", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()
	handler.Send(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Error reading response body:", err)
	}

	var errResp api.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatal("JSON unmarshal error:", err)
	}

	if errResp.Error != "Unsupported method!" {
		t.Errorf("unexpected error message: got %q, want %q", errResp.Error, "Unsupported method!")
	}
}
