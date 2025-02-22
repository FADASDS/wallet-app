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
	mux.HandleFunc("/api/send/", sendHandler.Send)
	mux.HandleFunc("/api/wallet/", walletHandler.ServeHTTP)
	mux.HandleFunc("/api/transactions/", transactionHandler.ServeHTTP)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
