package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
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
	log.Fatalln("[INFO] Creating postgres store")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	//Не стал отдельно добавлять проверки,
	//так как при отсутствии одной из переменных error всё равно выкинет sql.Open
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("[ERROR] Failed to connect postgres db: ", err)
	}
	store.db = db
	return &store, nil

}

func (s *storePG) Send(from string, to string, amount float64) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("[ERROR] Failed to open transaction: ", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("UPDATE wallet_tbl SET balance=balance-$1 WHERE id=$2", amount, from)
	if err != nil {
		log.Println("[ERROR] Failed to withdraw money: ", err)
		return err
	}

	_, err = tx.Exec("UPDATE wallet_tbl SET balance=balance+$1 WHERE id=$2", amount, to)
	if err != nil {
		log.Println("[ERROR] Failed to credit to an wallet: ", err)
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions_tbl (from_wallet, to_wallet, amount, date) VALUES ($1, $2, $3, $4)", from, to, amount, time.Now())
	if err != nil {
		log.Println("[ERROR] Failed to add transaction: ", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[ERROR] Failed to commit transaction: ", err)
		return err
	}

	return nil
}

func (s *storePG) GetLastNTransactions(n int64) ([]dto.TransactionDTO, error) {

	var transaction []dto.TransactionDTO
	rows, err := s.db.Query("SELECT * FROM transactions_tbl ORDER BY crtn_date LIMIT $1", n)

	if err != nil {
		log.Println("[ERROR] Failed to get transaction list: ", err)
		return nil, err
	}

	for rows.Next() {
		transactionItem := dto.TransactionDTO{}
		err = rows.Scan(&transactionItem.TransactionId, &transactionItem.From,
			&transactionItem.To, &transactionItem.Date, &transactionItem.Amount)
		if err != nil {
			log.Println("[ERROR] ОFailed to get transaction: ", err)
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
		log.Println("[ERROR] Failed to get wallet balance: ", err)
		return nil, err
	}

	return &walletBalance, nil
}
