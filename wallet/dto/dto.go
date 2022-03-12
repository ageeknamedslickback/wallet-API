package dto

import "github.com/shopspring/decimal"

// CrAmountInput is the credit amount input data transfer object
type CrAmountInput struct {
	CreditAmount decimal.Decimal `json:"credit_amount"`
}

// DrAmountInput is the debit amount input data transfer object
type DrAmountInput struct {
	DebitAmount decimal.Decimal `json:"debit_amount"`
}
