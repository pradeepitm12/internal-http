package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pradeepitm12/compaaa/internal-http/internal/domain/model"
	"github.com/shopspring/decimal"
)

type AccountRepo struct {
	db *sql.DB
}

// NewAccountRepository creates a new instance of the account repo.
func NewAccountRepository(db *sql.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) GetByID(ctx context.Context, tx *sql.Tx, id int) (*model.Account, error) {
	query := `SELECT account_id, balance FROM accounts WHERE account_id = $1 FOR UPDATE`

	var accID int
	var balance decimal.Decimal

	row := tx.QueryRowContext(ctx, query, id)
	if err := row.Scan(&accID, &balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	return model.NewAccount(accID, balance), nil
}

func (r *AccountRepo) Create(ctx context.Context, tx *sql.Tx, acc *model.Account) error {
	query := `INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`
	_, err := tx.ExecContext(ctx, query, acc.ID(), acc.Balance().String())
	if err != nil {
		fmt.Printf("error inserting account: %v\n", err)
		return err
	}

	return nil
}

func (r *AccountRepo) Update(ctx context.Context, tx *sql.Tx, acc *model.Account) error {
	query := `UPDATE accounts SET balance = $1 WHERE account_id = $2`
	_, err := tx.ExecContext(ctx, query, acc.Balance().String(), acc.ID())
	if err != nil {
		fmt.Println(err)
	}

	return err
}
