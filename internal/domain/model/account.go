package model

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Account struct {
	id      int
	balance decimal.Decimal
}

func NewAccount(id int, initial decimal.Decimal) *Account {
	return &Account{id: id, balance: initial}
}

func (a *Account) ID() int {
	return a.id
}

func (a *Account) Balance() decimal.Decimal {
	return a.balance
}

func (a *Account) Deposit(amount decimal.Decimal) {
	a.balance = a.balance.Add(amount)
}

func (a *Account) Withdraw(amount decimal.Decimal) error {
	if a.balance.LessThan(amount) {
		return fmt.Errorf("insufficient funds")
	}
	a.balance = a.balance.Sub(amount)
	return nil
}
