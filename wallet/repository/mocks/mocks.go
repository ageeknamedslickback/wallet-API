package mocks

import (
	"context"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/shopspring/decimal"
)

// MockRepo creates a mock the repository layer
type MockRepo struct {
	MockGetBalance func(
		ctx context.Context,
		walletID int,
	) (*domain.Wallet, error)
	MockUpdateBalance func(
		ctx context.Context,
		wallet *domain.Wallet,
		balance decimal.Decimal,
	) (*domain.Wallet, error)
}

// NewMockRepo inits a new instance of repository mocks with happy cases pre-defined
func NewMockRepo() *MockRepo {
	wallet := &domain.Wallet{
		ID:      1,
		Balance: decimal.NewFromFloat(200),
	}
	return &MockRepo{
		MockGetBalance: func(ctx context.Context, walletID int) (*domain.Wallet, error) { return wallet, nil },
		MockUpdateBalance: func(ctx context.Context, wallet *domain.Wallet, balance decimal.Decimal) (*domain.Wallet, error) {
			return wallet, nil
		},
	}
}

// GetBalance mocks GetBalance
func (m *MockRepo) GetBalance(
	ctx context.Context,
	walletID int,
) (*domain.Wallet, error) {
	return m.MockGetBalance(ctx, walletID)
}

// UpdateBalance mocks UpdateBalance
func (m *MockRepo) UpdateBalance(
	ctx context.Context,
	wallet *domain.Wallet,
	balance decimal.Decimal,
) (*domain.Wallet, error) {
	return m.MockUpdateBalance(ctx, wallet, balance)
}
