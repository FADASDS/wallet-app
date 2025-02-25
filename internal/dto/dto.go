// Package dto Пакет, предоставляющий структуры для ответов клиенту
package dto

// TransactionDTO структура для описания транзакции.
type TransactionDTO struct {
	TransactionId string  `json:"id"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Date          string  `json:"date"`
	Amount        float64 `json:"amount"`
}

// WalletBalanceDTO структура для описания ответа на запрос баланса кошелька.
type WalletBalanceDTO struct {
	WalletId string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}
