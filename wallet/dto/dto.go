package dto

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// AmountInput is the credit/debit amount input data transfer object
type AmountInput struct {
	Amount decimal.Decimal `json:"amount"`
}

// Valid validates the debit/credit amount is not a negative number
func (a *AmountInput) Valid() error {
	if a.Amount.IsNegative() {
		return fmt.Errorf("amount can not be a negative number")
	}
	return nil
}

// AccessToken represents Auth0 oauth2 access token
type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
