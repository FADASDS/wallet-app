package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"wallet-app/internal/lib/api"
	"wallet-app/internal/lib/api/response"
	"wallet-app/internal/storage"
)

type WalletHandler struct {
	Store storage.Storrer
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	path := req.URL.Path
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(path, "/")
	if len(parts) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(api.Error(http.StatusBadRequest, "No enough parameters (address)!"))
		w.Write(response)
		return
	}

	data, err := h.Store.GetWalletBalance(parts[3])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(api.Error(http.StatusNotFound, "GetWalletBalance error!"))
		w.Write(response)
		return
	}
	response, _ := json.Marshal(response.DataOkRsp(data))
	w.Write(response)

}

func (h *WalletHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodGet && strings.HasSuffix(req.URL.Path, "/balance"):
		h.GetBalance(w, req)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response, _ := json.Marshal(api.Error(http.StatusMethodNotAllowed, "Unsupported method!"))
		w.Write(response)
		return
	}
}
