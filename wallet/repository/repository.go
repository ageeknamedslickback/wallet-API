package repository

import (
	"context"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/shopspring/decimal"
)

// Get represents a contract for all GET operations in the infra database layer
type Get interface {
	GetBalance(
		ctx context.Context,
		walletID int,
	) (*domain.Wallet, error)
}

// Update represents a contract for all UPDATE operations in the infra database layer
type Update interface {
	UpdateBalance(
		ctx context.Context,
		wallet *domain.Wallet,
		balance decimal.Decimal,
	) (*domain.Wallet, error)
}
