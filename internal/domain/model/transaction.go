package model

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionID int64

type Transaction struct {
	transactionID   int64
	sourceAccountID int
	destAccountID   int
	amount          decimal.Decimal
	createdAt       time.Time
}

func NewTransaction(sourceID, destID int, amount decimal.Decimal) (*Transaction, error) {
	if sourceID == destID {
		return nil, errors.New("source and destination accounts must differ")
	}
	return &Transaction{
		sourceAccountID: sourceID,
		destAccountID:   destID,
		amount:          amount,
		createdAt:       time.Now().UTC(),
	}, nil
}

func (t *Transaction) ID() int64 {
	return t.transactionID
}

func (t *Transaction) SourceAccountID() int {
	return t.sourceAccountID
}

func (t *Transaction) DestAccountID() int {
	return t.destAccountID
}

func (t *Transaction) Amount() decimal.Decimal {
	return t.amount
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Transaction) SetID(id int64) {
	t.transactionID = id
}
