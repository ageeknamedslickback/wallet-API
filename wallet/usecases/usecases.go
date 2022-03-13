package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/ageeknamedslickback/wallet-API/wallet/repository"
	"github.com/shopspring/decimal"
)

// WalletBusinessLogic designs wallet's business logic that has been implemented
type WalletBusinessLogic interface {
	WalletBalance(
		ctx context.Context,
		walletID int,
	) (*decimal.Decimal, error)
	CreditWallet(
		ctx context.Context,
		walletID int,
		creditAmount decimal.Decimal,
	) (*domain.Wallet, error)
	DebitWallet(
		ctx context.Context,
		walletID int,
		debitAmount decimal.Decimal,
	) (*domain.Wallet, error)
}

// WalletUsecases sets up wallet's API server usecase layer
// with all the necessary dependencies
type WalletUsecases struct {
	Get    repository.Get
	Update repository.Update
}

// NewWalletUsecases initializes wallet's business logic
func NewWalletUsecases(
	get repository.Get,
	update repository.Update,
) *WalletUsecases {
	w := &WalletUsecases{
		Get:    get,
		Update: update,
	}
	w.checkPreconditions()
	return w
}

func (w *WalletUsecases) checkPreconditions() {
	if w.Get == nil {
		log.Panicf("wallet usecases have not initalized GET repository")
	}
	if w.Update == nil {
		log.Panicf("wallet usecases have not initalized UPDATE repository")
	}
}

// WalletBalance gets the current balance of a wallet
func (w *WalletUsecases) WalletBalance(
	ctx context.Context,
	walletID int,
) (*decimal.Decimal, error) {
	wallet, err := w.Get.GetWallet(ctx, walletID)
	if err != nil {
		return nil, err
	}

	return &wallet.Balance, nil
}

// CreditWallet credits money on a given wallet
func (w *WalletUsecases) CreditWallet(
	ctx context.Context,
	walletID int,
	creditAmount decimal.Decimal,
) (*domain.Wallet, error) {
	wallet, err := w.Get.GetWallet(ctx, walletID)
	if err != nil {
		return nil, err
	}
	balance := wallet.Balance.Sub(creditAmount)

	if balance.IsNegative() {
		return nil, fmt.Errorf("a wallet balance cannot go below 0")
	}

	return w.Update.UpdateBalance(ctx, wallet, balance)
}

// DebitWallet debits money on a given wallet
func (w *WalletUsecases) DebitWallet(
	ctx context.Context,
	walletID int,
	debitAmount decimal.Decimal,
) (*domain.Wallet, error) {
	wallet, err := w.Get.GetWallet(ctx, walletID)
	if err != nil {
		return nil, err
	}
	balance := wallet.Balance.Add(debitAmount)

	return w.Update.UpdateBalance(ctx, wallet, balance)
}
