package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wallet-app/internal/lib/api"
	"wallet-app/internal/lib/api/response"
	"wallet-app/internal/storage"
)

type TransactionHandler struct {
	Store storage.Storrer
}

func (t *TransactionHandler) GetLastNTransactions(w http.ResponseWriter, req *http.Request) {
	count := req.URL.Query().Get("count")
	if count == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(api.Error(http.StatusBadRequest, "count parameter is missing"))
		w.Write(response)
		return
	}
	n, err := strconv.ParseInt(count, 10, 64)

	if err != nil || n <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(api.Error(http.StatusBadRequest, "Invalid count parameter"))
		w.Write(response)
		return
	}

	data, err := t.Store.GetLastNTransactions(n)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(api.Error(http.StatusNotFound, "GetLastNTransactions error!"))
		w.Write(response)
		return
	}
	response, _ := json.Marshal(response.DataOkRsp(data))
	w.Write(response)
}

func (t *TransactionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch {
	case req.Method == http.MethodGet:
		t.GetLastNTransactions(w, req)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response, _ := json.Marshal(api.Error(http.StatusMethodNotAllowed, "Unsupported method!"))
		w.Write(response)
		return
	}
}
