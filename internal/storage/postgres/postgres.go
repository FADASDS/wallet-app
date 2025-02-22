package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"wallet-app/internal/dto"
)

// TODO перейти с прямых запросов на миграции
type storePG struct {
	db *sql.DB
}

func NewStorage() (*storePG, error) {
	var store storePG

	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	store.db = db
	return &store, nil

}

func (s *storePG) Send(from string, to string, amount float64) error {

	return nil
}

func (s *storePG) GetLastNTransactions(n int64) ([]dto.TransactionDTO, error) {

	var transaction []dto.TransactionDTO
	rows, err := s.db.Query("SELECT * FROM wallet_db ORDER BY  crtn_date LIMIT $1", n)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		transactionItem := dto.TransactionDTO{}
		err = rows.Scan(&transactionItem.TransactionId, &transactionItem.From,
			&transactionItem.To, &transactionItem.Date, &transactionItem.Amount)
		if err != nil {
			return nil, err
		}
		transaction = append(transaction, transactionItem)
	}

	return transaction, nil
}

func (s *storePG) GetWalletBalance(address string) (*dto.WalletBalanceDTO, error) {
	var walletBalance dto.WalletBalanceDTO
	row := s.db.QueryRow("SELECT id, balance FROM wallet_tbl WHERE id = $1", address)

	err := row.Scan(&walletBalance.WalletId, &walletBalance.Balance)

	if err != nil {
		return nil, err
	}

	return &walletBalance, nil
}
