package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"wallet-app/internal/lib/api"
	"wallet-app/internal/lib/api/response"
	"wallet-app/internal/storage"
)

const (
	msgErr = "Wrong json"
)

type SendHandler struct {
	Store storage.Storrer
}

type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func (h *SendHandler) Send(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		log.Println("[INFO] Received unsupported request")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response, _ := json.Marshal(api.Error(http.StatusMethodNotAllowed, "Unsupported method!"))
		w.Write(response)
		return
	}

	log.Println("[INFO] Received send request")

	w.Header().Set("Content-Type", "application/json")

	defer req.Body.Close()
	var s Transaction
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("[ERROR] Failed to read request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		response, _ := json.Marshal(api.Error(http.StatusInternalServerError, "Internal server error."))
		w.Write(response)
		return
	}

	err = json.Unmarshal(body, &s)

	if err != nil {
		log.Println("[ERROR] Failed to Unmarshal request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		response, _ := json.Marshal(api.Error(http.StatusInternalServerError, "Internal server error."))
		w.Write(response)
		return
	}

	if s.Amount <= 0 {
		log.Println("[ERROR] Invalid amount for send retrieved")
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(api.Error(http.StatusBadRequest, "Invalid amount"))
		w.Write(response)
		return
	}

	err = h.Store.Send(s.From, s.To, s.Amount)
	if err != nil {
		log.Println("[ERROR] Failed to send money: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		response, _ := json.Marshal(api.Error(http.StatusInternalServerError, "Internal server error."))
		w.Write(response)
		return
	}

	log.Println("[INFO] Money sent successfully.")

	response, _ := json.Marshal(response.EmptyOkRsp())
	w.Write(response)
}
