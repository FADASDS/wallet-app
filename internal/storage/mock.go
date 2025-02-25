// Package storage содержит структуры и методы для работы с различными хранилищами данных.
package storage

import (
	"fmt"
	"wallet-app/internal/dto"
)

// MockStore структура-заглушка для тестирования хэндлеров.
type MockStore struct{}

// GetLastNTransactions имитирует функцию GetLastNTransactions для тестирования хэндлеров.
func (m *MockStore) GetLastNTransactions(n int64) ([]dto.TransactionDTO, error) {
	return []dto.TransactionDTO{
		{TransactionId: "tx-1", From: "wallet1", To: "wallet2", Amount: 100.0},
		{TransactionId: "tx-2", From: "wallet3", To: "wallet4", Amount: 50.0},
		{TransactionId: "tx-3", From: "wallet5", To: "wallet6", Amount: 150.0},
	}, nil
}

// Send имитирует функцию Send для тестирования хэндлеров.
func (m *MockStore) Send(from string, to string, amount float64) error {
	walletId1 := "wallet_1"
	walletId2 := "wallet_2"
	balance1 := 100.0

	if from != walletId1 {
		return fmt.Errorf("from is not wallet1")
	}
	if to != walletId2 {
		return fmt.Errorf("to is not wallet2")
	}
	if amount > balance1 {
		return fmt.Errorf("amount > balance1")
	}
	return nil
}

// GetWalletBalance имитирует функцию GetWalletBalance для тестирования хэндлеров.
func (m *MockStore) GetWalletBalance(address string) (*dto.WalletBalanceDTO, error) {
	if address == "wallet_1" {
		return &dto.WalletBalanceDTO{WalletId: "wallet_1", Balance: 100.0}, nil
	}
	return nil, fmt.Errorf("Wrong wallet address")
}
