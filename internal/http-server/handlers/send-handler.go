package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"wallet-app/internal/lib/api"
	"wallet-app/internal/storage"
)

const (
	msgErr = "Wrong json"
)

type SendHandler struct {
	Store storage.Storrer
}

type Transaction struct {
	from   string  `json:"from"`
	to     string  `json:"to"`
	amount float64 `json:"amount"`
}

func (h *SendHandler) Send(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	defer req.Body.Close()
	var s Transaction
	body, err := io.ReadAll(req.Body)
	if err != nil {
		response, _ := json.Marshal(api.Error(http.StatusInternalServerError, "Internal server error"))
		w.Write(response)
		return
	}
	err = json.Unmarshal(body, &s)
	if err != nil {
		response, _ := json.Marshal(api.Error(http.StatusInternalServerError, "Internal server error"))
		w.Write(response)
		return
	}

	h.Store.Send(s.from, s.to, s.amount)
}
