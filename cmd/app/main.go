package main

import (
	"log"
	"net/http"
	"wallet-app/internal/http-server/handlers"
	"wallet-app/internal/storage/postgres"
)

func main() {
	store, err := postgres.NewStorage()
	walletHandler := &handlers.WalletHandler{Store: store}
	sendHandler := &handlers.SendHandler{Store: store}
	transactionHandler := &handlers.TransactionHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/send", sendHandler.Send)
	mux.HandleFunc("/api/wallet/", walletHandler.ServeHTTP)
	mux.HandleFunc("/api/transactions", transactionHandler.ServeHTTP)

	log.Println("[INFO] Starting web-server http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
