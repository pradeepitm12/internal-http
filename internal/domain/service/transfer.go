package service

import (
	"context"

	"github.com/pradeepitm12/compaaa/internal-http/internal/domain/model"
	"github.com/shopspring/decimal"
)

type TransferService interface {
	Transfer(ctx context.Context, sourceID, destID int, amount decimal.Decimal) (*model.Transaction, error)
}
