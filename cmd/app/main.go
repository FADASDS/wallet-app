// Package main является точкой входа в приложение
package main

import (
	"log"
	"net/http"
	"wallet-app/internal/http-server/handlers"
	"wallet-app/internal/storage/postgres"
)

// Функция, являющаяся точкой входа в приложение. Запускает сервер, создаёт экземпляр хранилища и обработчики
func main() {
	store, err := postgres.NewStorage()
	if err != nil {
		log.Fatal("[ERROR] Failed to connect: ", err)
	}
	defer store.Db.Close()

	walletHandler := &handlers.WalletHandler{Store: store}
	sendHandler := &handlers.SendHandler{Store: store}
	transactionHandler := &handlers.TransactionHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/send", sendHandler.Send)
	mux.HandleFunc("/api/wallet/", walletHandler.ServeHTTP)
	mux.HandleFunc("/api/transactions", transactionHandler.ServeHTTP)

	err = http.ListenAndServe(":8080", mux)
	log.Println("[INFO] Starting web-server http://localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
}
