package tests

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
	"wallet-app/internal/storage/postgres"
)

func TestGetBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := &postgres.StorePG{Db: db}
	address := "wallet_1"
	rows := sqlmock.NewRows([]string{"id", "balance"}).AddRow(address, 100.0)
	mock.ExpectQuery(`SELECT id, balance FROM wallet_tbl WHERE id \= \$1`).
		WithArgs(address).
		WillReturnRows(rows)

	res, err := store.GetWalletBalance(address)

	if err != nil {
		t.Fatal(err)
	}

	if res == nil {
		t.Fatal("balance is nil")
	}

	if address != res.WalletId {
		t.Fatalf("Wrong address, expected %s, got %s", address, res.WalletId)
	}

	if res.Balance != 100.0 {
		t.Fatalf("Wrong balance, expected %f, got %f", 100.0, res.Balance)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBalanceWrongAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := &postgres.StorePG{Db: db}
	address := "wrong_wallet"
	rows := sqlmock.NewRows([]string{"id", "balance"})
	mock.ExpectQuery(`SELECT id, balance FROM wallet_tbl WHERE id \= \$1`).
		WithArgs(address).
		WillReturnRows(rows)

	res, err := store.GetWalletBalance(address)

	if err != nil {
		t.Fatal(err)
	}

	if res != nil {
		t.Fatalf("balance must be nil, got %+v", res)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLastNTransactions(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := &postgres.StorePG{Db: db}
	count := int64(2)
	rows := sqlmock.NewRows([]string{"id", "from", "to", "crtn_date", "amount"}).
		AddRow("tx_1", "wallet_1", "wallet_2", time.Now(), 100.0).
		AddRow("tx_2", "wallet_3", "wallet_4", time.Now(), 100.0)

	mock.ExpectQuery(`SELECT \* FROM transactions_tbl ORDER BY crtn_date LIMIT \$1`).
		WithArgs(count).
		WillReturnRows(rows)

	res, err := store.GetLastNTransactions(count)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	if len(res) != 2 {
		t.Fatalf("Not enough transactions, expected %d, got: %d", count, len(res))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLastNTransactionsWrongCount(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := &postgres.StorePG{Db: db}
	count := int64(3)
	rows := sqlmock.NewRows([]string{"id", "from", "to", "crtn_date", "amount"})

	mock.ExpectQuery(`SELECT \* FROM transactions_tbl ORDER BY crtn_date LIMIT \$1`).
		WithArgs(count).
		WillReturnRows(rows)

	res, err := store.GetLastNTransactions(count)

	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 0 {
		t.Fatalf("Not enough transactions, expected %d, got: %d", len(res), count)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSend(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	from := "wallet_1"
	to := "wallet_2"
	amount := 50.0
	store := &postgres.StorePG{Db: db}
	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE wallet_tbl SET balance=balance\-\$1 WHERE id=\$2`).
		WithArgs(amount, from).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE wallet_tbl SET balance=balance\+\$1 WHERE id=\$2`).
		WithArgs(amount, to).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO transactions_tbl`).
		WithArgs(from, to, amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = store.Send(from, to, amount)
	if err != nil {
		t.Fatal(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendRollback(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	from := "wallet_1"
	to := "wallet_2"
	amount := 50.0
	store := &postgres.StorePG{Db: db}
	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE wallet_tbl SET balance=balance\-\$1 WHERE id=\$2`).
		WithArgs(amount, from).
		WillReturnError(fmt.Errorf("Failed to withdraw"))

	mock.ExpectRollback()

	err = store.Send(from, to, amount)
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}
