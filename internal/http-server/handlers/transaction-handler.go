package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"wallet-app/internal/lib/api/response"
	"wallet-app/internal/storage"
)

// TransactionHandler cтруктура для описания обработчика для запроса на получение N последних транзакций.
type TransactionHandler struct {
	Store storage.Storrer
}

// GetLastNTransactions метод отвечающий за обработку запроса.
func (t *TransactionHandler) GetLastNTransactions(w http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] Received get last n transactions request")

	count := req.URL.Query().Get("count")
	if count == "" {
		log.Println("[ERROR] Failed to get 'count' parameter")
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(response.Error(http.StatusBadRequest, "'count' parameter is missing."))
		w.Write(response)
		return
	}
	n, err := strconv.ParseInt(count, 10, 64)

	if err != nil || n <= 0 {
		if n > 0 {
			log.Println("[ERROR] Failed to parse 'count' parameter", err)
		} else {
			log.Println("[ERROR] Invalid 'count' parameter")
		}
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(response.Error(http.StatusBadRequest, "Invalid 'count' parameter."))
		w.Write(response)
		return
	}

	data, err := t.Store.GetLastNTransactions(n)
	if err != nil {
		log.Println("[ERROR] Failed to get last transactions", err)
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(response.Error(http.StatusNotFound, "Get last NTransactions error!"))
		w.Write(response)
		return
	}

	log.Println("[INFO] Transaction success.")
	response, _ := json.Marshal(response.DataOkRsp(data))
	w.Write(response)
}

// ServeHTTP метод отвечающий за маршрутизацию запроса.
func (t *TransactionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch {
	case req.Method == http.MethodGet:
		t.GetLastNTransactions(w, req)
		return

	default:
		log.Println("[INFO] Received unsupported request")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response, _ := json.Marshal(response.Error(http.StatusMethodNotAllowed, "Unsupported method!"))
		w.Write(response)
		return
	}
}
