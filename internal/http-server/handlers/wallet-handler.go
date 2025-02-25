package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"wallet-app/internal/lib/api/response"
	"wallet-app/internal/storage"
)

// WalletHandler cтруктура для описания обработчика для запроса на получение баланса кошелька.
type WalletHandler struct {
	Store storage.Storrer
}

// GetBalance метод отвечающий за обработку запроса.
func (h *WalletHandler) GetBalance(w http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] Received get wallet balance request")

	w.Header().Set("Content-Type", "application/json")

	address := strings.TrimPrefix(req.URL.Path, "/api/wallet/")
	address = strings.TrimSuffix(address, "/balance")

	if address == "" {
		log.Println("[ERROR] Invalid address")
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(response.Error(http.StatusBadRequest, "Not enough parameters (address)!"))
		w.Write(response)
		return
	}
	log.Println("[INFO] Getting wallet balance...")
	data, err := h.Store.GetWalletBalance(address)
	if err != nil {
		log.Println("[ERROR] Failed to get balance: ", err)
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(response.Error(http.StatusNotFound, "GetWalletBalance error!"))
		w.Write(response)
		return
	}

	log.Println("[INFO] Wallet balance retrieved successfully.")
	response, _ := json.Marshal(response.DataOkRsp(data))
	w.Write(response)

}

// ServeHTTP метод отвечающий за маршрутизацию запроса.
func (h *WalletHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodGet && strings.HasSuffix(req.URL.Path, "/balance"):
		h.GetBalance(w, req)
		return
	default:
		log.Println("[INFO] Received unsupported request")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response, _ := json.Marshal(response.Error(http.StatusMethodNotAllowed, "Unsupported method!"))
		w.Write(response)
		return
	}
}
