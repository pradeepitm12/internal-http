package util

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func StringToDecimal(s string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("invalid decimal string: %w", err)
	}
	return d, nil
}
