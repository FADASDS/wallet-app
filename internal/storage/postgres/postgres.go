package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
	"wallet-app/internal/dto"
)

// TODO перейти с прямых запросов на миграции
type storePG struct {
	db *sql.DB
}

func NewStorage() (*storePG, error) {
	var store storePG

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	//Не стал отдельно добавлять проверки,
	//так как при отсутствии одной из переменных error всё равно выкинет sql.Open
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	store.db = db
	return &store, nil

}

func (s *storePG) Send(from string, to string, amount float64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("UPDATE wallet_tbl SET balance=balance-$1 WHERE id=$2", amount, from)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE wallet_tbl SET balance=balance+$1 WHERE id=$2", amount, to)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions_tbl (from_wallet, to_wallet, amount, date) VALUES ($1, $2, $3, $4)", from, to, amount, time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *storePG) GetLastNTransactions(n int64) ([]dto.TransactionDTO, error) {

	var transaction []dto.TransactionDTO
	rows, err := s.db.Query("SELECT * FROM wallet_tbl ORDER BY crtn_date LIMIT $1", n)

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
