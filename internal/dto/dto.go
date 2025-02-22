package dto

type TransactionDTO struct {
	TransactionId string  `json:"id"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Date          string  `json:"date"`
	Amount        float64 `json:"amount"`
}

type WalletBalanceDTO struct {
	WalletId string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}
