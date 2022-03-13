package usecases_test

import (
	"context"
	"log"
	"testing"

	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/database"
	"github.com/ageeknamedslickback/wallet-API/wallet/usecases"
	"github.com/shopspring/decimal"
)

func initTestUsecases() *usecases.WalletUsecases {
	gormDb, err := database.ConnectToDatabase()
	if err != nil {
		log.Panicf("error connecting to the database: %v", err)
	}
	getRepo := database.NewWalletDb(gormDb)
	updateRepo := database.NewWalletDb(gormDb)
	w := usecases.NewWalletUsecases(getRepo, updateRepo)
	return w
}

func TestWalletUsecases_WalletBalance(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		walletID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:      ctx,
				walletID: 1, // exists in the database
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:      ctx,
				walletID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := initTestUsecases()
			balance, err := w.WalletBalance(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletUsecases.WalletBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if balance == nil {
					t.Fatalf("expected a wallet")
				}

				exectedBal := decimal.NewFromInt(100).String()
				if balance.String() != exectedBal {
					t.Fatalf(
						"expected wallet balance to be 100 but got %s",
						balance,
					)
				}
			}

			if tt.wantErr {
				if balance != nil {
					t.Fatalf("expected no wallet balance to be returned")
				}
			}
		})
	}
}

func TestWalletUsecases_CreditWallet(t *testing.T) {
	w := initTestUsecases()
	ctx := context.Background()
	amount := decimal.NewFromFloat(10.34)
	walletID := 2

	balance, err := w.WalletBalance(ctx, walletID)
	if err != nil {
		t.Fatalf("failed to get wallet balance with id 2: %v", err)
	}

	type args struct {
		ctx          context.Context
		walletID     int
		creditAmount decimal.Decimal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:          ctx,
				walletID:     walletID,
				creditAmount: amount,
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:          ctx,
				walletID:     0,
				creditAmount: amount,
			},
			wantErr: true,
		},
		{
			name: "sad case - negative amount",
			args: args{
				ctx:          ctx,
				walletID:     0,
				creditAmount: decimal.NewFromFloat(-10.34),
			},
			wantErr: true,
		},
		{
			name: "sad case - balance cannot go below 0",
			args: args{
				ctx:          ctx,
				walletID:     0,
				creditAmount: decimal.NewFromFloat(1034),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := w.CreditWallet(tt.args.ctx, tt.args.walletID, tt.args.creditAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletUsecases.CreditWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				expectedBalance := balance.Sub(amount)
				if wallet.Balance.String() != expectedBalance.String() {
					t.Fatalf("expected amount to be credited to the wallet balance")
				}

				// restore the amount
				if _, err := w.DebitWallet(ctx, walletID, amount); err != nil {
					t.Fatalf("error restoring the credited amount: %v", err)
				}
			}

			if tt.wantErr && wallet != nil {
				t.Fatalf("did not expect a wallet")
			}
		})
	}
}

func TestWalletUsecases_DebitWallet(t *testing.T) {
	w := initTestUsecases()
	ctx := context.Background()
	amount := decimal.NewFromFloat(50)
	walletID := 2

	balance, err := w.WalletBalance(ctx, walletID)
	if err != nil {
		t.Fatalf("failed to get wallet balance with id 2: %v", err)
	}

	type args struct {
		ctx         context.Context
		walletID    int
		debitAmount decimal.Decimal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:         ctx,
				walletID:    walletID,
				debitAmount: amount,
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:         ctx,
				walletID:    0,
				debitAmount: amount,
			},
			wantErr: true,
		},
		{
			name: "sad case - negative amount",
			args: args{
				ctx:         ctx,
				walletID:    0,
				debitAmount: decimal.NewFromFloat(-10.34),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := w.DebitWallet(tt.args.ctx, tt.args.walletID, tt.args.debitAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletUsecases.DebitWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				expectedBalance := balance.Add(amount)
				if wallet.Balance.String() != expectedBalance.String() {
					t.Fatalf("expected amount to be debited to the wallet balance")
				}

				// restore the amount
				if _, err := w.CreditWallet(ctx, walletID, amount); err != nil {
					t.Fatalf("error restoring the debited amount: %v", err)
				}
			}

			if tt.wantErr && wallet != nil {
				t.Fatalf("did not expect a wallet")
			}
		})
	}
}
