package types

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type NumberStr string

func (t NumberStr) MarshalJSON() ([]byte, error) {
	tt, _ := decimal.NewFromString(string(t))
	s := fmt.Sprintf("\"%s\"", tt.String())
	return []byte(s), nil
}
