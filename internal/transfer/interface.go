package transfer

import (
	"context"
	"database/sql"

	"github.com/pradeepitm12/compaaa/internal-http/internal/domain/model"
)

type AccountRepository interface {
	GetByID(ctx context.Context, tx *sql.Tx, id int) (*model.Account, error)
	Update(ctx context.Context, tx *sql.Tx, account *model.Account) error
	Create(ctx context.Context, tx *sql.Tx, account *model.Account) error
}

type TransactionRepository interface {
	Save(ctx context.Context, tx *model.Transaction) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}
