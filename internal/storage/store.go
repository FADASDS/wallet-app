package storage

import "wallet-app/internal/dto"

type Storrer interface {
	Send(from string, to string, amount float64) error
	GetWalletBalance(address string) (*dto.WalletBalanceDTO, error)

	GetLastNTransactions(n int64) ([]dto.TransactionDTO, error)
}
