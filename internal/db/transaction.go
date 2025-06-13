package db

import (
	"context"
	"database/sql"

	"github.com/pradeepitm12/compaaa/internal-http/internal/domain/model"
)

type TransactionRepo struct {
	db *sql.DB
}

// NewTransactionRepository creates a new instance of the transaction repo.
func NewTransactionRepository(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Save(ctx context.Context, tx *model.Transaction) error {
	query := `
		INSERT INTO transactions (source_account_id, destination_account_id, amount, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING transaction_id
	`

	row := r.db.QueryRowContext(ctx, query, tx.SourceAccountID(), tx.DestAccountID(), tx.Amount(), tx.CreatedAt())

	var id int64
	if err := row.Scan(&id); err != nil {
		return err
	}

	tx.SetID(id)
	return nil
}
