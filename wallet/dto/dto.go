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

// WrappedError is a custom context wrapped error
type WrappedError struct {
	Context string `json:"context"`
	Err     error  `json:"error"`
}

// Error is a string representation of an error interface
func (w *WrappedError) Error() string {
	return fmt.Sprintf("%s: %v", w.Context, w.Err)
}

// Wrap wraps an error with it's context
func Wrap(err error, info string) *WrappedError {
	return &WrappedError{
		Context: info,
		Err:     err,
	}
}
