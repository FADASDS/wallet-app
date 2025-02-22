package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"wallet-app/internal/lib/api"
	"wallet-app/internal/storage"
)

type TransactionHandler struct {
	Store storage.Storrer
}

func (t *TransactionHandler) GetLastNTransactions(w http.ResponseWriter, req *http.Request) {
	count := req.URL.Query().Get("count")
	if count == "" {
		response, _ := json.Marshal(api.Error(http.StatusBadRequest, "count parameter is missing"))
		w.Write(response)
		return
	}
	n, err := strconv.ParseInt(count, 10, 64)

	if err != nil {
		response, _ := json.Marshal(api.Error(http.StatusInternalServerError, "Internal server error"))
		w.Write(response)
	}

	t.Store.GetLastNTransactions(n)
}

func (t *TransactionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	w.Header().Set("Content-Type", "application/json")
	switch {
	case req.Method == http.MethodGet && strings.HasSuffix(path, "/balance"):
		parts := strings.Split(path, "/")
		if len(parts) != 4 {
			response, _ := json.Marshal(api.Error(http.StatusBadRequest, "No enough parameters!"))
			w.Write(response)
			return
		}
		t.GetLastNTransactions(w, req)
		return

	default:
		response, _ := json.Marshal(api.Error(http.StatusMethodNotAllowed, "Unsupported method!"))
		w.Write(response)
		return
	}
}
