package domain

import (
	"github.com/shopspring/decimal"
)

// Wallet represents a digital wallet that manages
// debit and credit transaction for online casino game players
type Wallet struct {
	ID      uint            `json:"id" gorm:"primarykey"`
	Balance decimal.Decimal `json:"balance"`
}
