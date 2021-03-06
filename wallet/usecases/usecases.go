package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/ageeknamedslickback/wallet-API/wallet/dto"
	"github.com/ageeknamedslickback/wallet-API/wallet/repository"
	"github.com/shopspring/decimal"
)

// WalletBusinessLogic designs wallet's business logic that has been implemented
type WalletBusinessLogic interface {
	WalletBalance(
		ctx context.Context,
		walletID int,
	) (*domain.Wallet, error)
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
) (*domain.Wallet, error) {
	wallet, err := w.Get.GetBalance(ctx, walletID)
	if err != nil {
		return nil, dto.Wrap(err, "WalletBalance")
	}

	return wallet, nil
}

// CreditWallet credits money on a given wallet
func (w *WalletUsecases) CreditWallet(
	ctx context.Context,
	walletID int,
	creditAmount decimal.Decimal,
) (*domain.Wallet, error) {
	wallet, err := w.Get.GetBalance(ctx, walletID)
	if err != nil {
		return nil, dto.Wrap(err, "CreditWallet")
	}
	balance := wallet.Balance.Sub(creditAmount)

	if balance.IsNegative() {
		return nil, dto.Wrap(fmt.Errorf("a wallet balance cannot go below 0"), "CreditWallet")
	}

	updatedWallet, err := w.Update.UpdateBalance(ctx, wallet, balance)
	if err != nil {
		return nil, dto.Wrap(err, "CreditWallet")
	}

	return updatedWallet, nil
}

// DebitWallet debits money on a given wallet
func (w *WalletUsecases) DebitWallet(
	ctx context.Context,
	walletID int,
	debitAmount decimal.Decimal,
) (*domain.Wallet, error) {
	wallet, err := w.Get.GetBalance(ctx, walletID)
	if err != nil {
		return nil, dto.Wrap(err, "DebitWallet")
	}
	balance := wallet.Balance.Add(debitAmount)

	updatedWallet, err := w.Update.UpdateBalance(ctx, wallet, balance)
	if err != nil {
		dto.Wrap(err, "DebitWallet")
	}

	return updatedWallet, nil
}
