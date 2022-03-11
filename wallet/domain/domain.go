package domain

import (
	"gorm.io/gorm"
)

// Wallet represents a digital wallet that manages
// debit and credit transaction for online casino game players
type Wallet struct {
	gorm.Model
}
