package usecases_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/database"
	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/services/cache"
	"github.com/ageeknamedslickback/wallet-API/wallet/repository/mocks"
	"github.com/ageeknamedslickback/wallet-API/wallet/usecases"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
)

var ctx = context.Background()

func initTestUsecases() *usecases.WalletUsecases {
	gormDb, err := database.ConnectToDatabase()
	if err != nil {
		log.Panicf("error connecting to the database: %v", err)
	}

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	c := cache.NewCacheService(rdb)
	getRepo := database.NewWalletDb(gormDb, c)
	updateRepo := database.NewWalletDb(gormDb, c)
	w := usecases.NewWalletUsecases(getRepo, updateRepo)
	return w
}

func UnitTestWalletBalance(t *testing.T) {
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
			getMockRepo := mocks.NewMockRepo()
			updateMockRepo := mocks.NewMockRepo()
			w := usecases.NewWalletUsecases(getMockRepo, updateMockRepo)

			if tt.name == "sad case" {
				getMockRepo.MockGetBalance = func(ctx context.Context, walletID int) (*domain.Wallet, error) {
					return nil, fmt.Errorf("failed to get wallet")
				}
			}

			wallet, err := w.WalletBalance(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletUsecases.WalletBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if wallet == nil {
					t.Fatalf("expected a wallet")
				}

				exectedBal := decimal.NewFromInt(200).String()
				if wallet.Balance.String() != exectedBal {
					t.Fatalf(
						"expected wallet balance to be 200 but got %s",
						wallet.Balance,
					)
				}
			}

			if tt.wantErr {
				if wallet != nil {
					t.Fatalf("expected no wallet balance to be returned")
				}
			}
		})
	}
}

func TestWalletUsecases_WalletBalance(t *testing.T) {
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
			wallet, err := w.WalletBalance(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletUsecases.WalletBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if wallet == nil {
					t.Fatalf("expected a wallet")
				}

				exectedBal := decimal.NewFromInt(100).String()
				if wallet.Balance.String() != exectedBal {
					t.Fatalf(
						"expected wallet balance to be 100 but got %s",
						wallet.Balance,
					)
				}
			}

			if tt.wantErr {
				if wallet != nil {
					t.Fatalf("expected no wallet balance to be returned")
				}
			}
		})
	}
}

func UnitTestCreditWallet(t *testing.T) {
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
				walletID:     1,
				creditAmount: decimal.NewFromFloat(100),
			},
			wantErr: false,
		},
		{
			name: "sad case; failed to get balance",
			args: args{
				ctx:          ctx,
				walletID:     1,
				creditAmount: decimal.NewFromFloat(100),
			},
			wantErr: true,
		},
		{
			name: "sad case - failed to update balance",
			args: args{
				ctx:          ctx,
				walletID:     1,
				creditAmount: decimal.NewFromFloat(10.34),
			},
			wantErr: true,
		},
		{
			name: "sad case - balance below 0",
			args: args{
				ctx:          ctx,
				walletID:     1,
				creditAmount: decimal.NewFromFloat(1000.34),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockRepo := mocks.NewMockRepo()
			updateMockRepo := mocks.NewMockRepo()
			w := usecases.NewWalletUsecases(getMockRepo, updateMockRepo)

			if tt.name == "happy case" {
				updateMockRepo.MockUpdateBalance = func(ctx context.Context, wallet *domain.Wallet, balance decimal.Decimal) (*domain.Wallet, error) {
					return &domain.Wallet{ID: 1, Balance: decimal.NewFromFloat(250)}, nil
				}
			}

			if tt.name == "sad case; failed to get balance" {
				getMockRepo.MockGetBalance = func(ctx context.Context, walletID int) (*domain.Wallet, error) { return nil, fmt.Errorf("error") }
			}

			if tt.name == "sad case - failed to update balance" {
				updateMockRepo.MockUpdateBalance = func(ctx context.Context, wallet *domain.Wallet, balance decimal.Decimal) (*domain.Wallet, error) {
					return nil, fmt.Errorf("error")
				}
			}

			wallet, err := w.CreditWallet(tt.args.ctx, tt.args.walletID, tt.args.creditAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletUsecases.CreditWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if wallet == nil {
					t.Fatalf("expected a wallet")
				}

				if wallet.Balance.String() != decimal.NewFromFloat(250).String() {
					t.Fatalf("wallet balance was not updated")
				}
			}

			if tt.wantErr && wallet != nil {
				t.Fatalf("did not expect a wallet")
			}
		})
	}
}

func TestWalletUsecases_CreditWallet(t *testing.T) {
	w := initTestUsecases()
	amount := decimal.NewFromFloat(10.34)
	walletID := 2

	wal, err := w.WalletBalance(ctx, walletID)
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
				expectedBalance := wal.Balance.Sub(amount)
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

func UnitTestDebitWallet(t *testing.T) {
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
				walletID:    1,
				debitAmount: decimal.NewFromFloat(100),
			},
			wantErr: false,
		},
		{
			name: "sad case; failed to get balance",
			args: args{
				ctx:         ctx,
				walletID:    1,
				debitAmount: decimal.NewFromFloat(100),
			},
			wantErr: true,
		},
		{
			name: "sad case - failed to update balance",
			args: args{
				ctx:         ctx,
				walletID:    1,
				debitAmount: decimal.NewFromFloat(10.34),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockRepo := mocks.NewMockRepo()
			updateMockRepo := mocks.NewMockRepo()
			w := usecases.NewWalletUsecases(getMockRepo, updateMockRepo)

			if tt.name == "happy case" {
				updateMockRepo.MockUpdateBalance = func(ctx context.Context, wallet *domain.Wallet, balance decimal.Decimal) (*domain.Wallet, error) {
					return &domain.Wallet{ID: 1, Balance: decimal.NewFromFloat(250)}, nil
				}
			}

			if tt.name == "sad case; failed to get balance" {
				getMockRepo.MockGetBalance = func(ctx context.Context, walletID int) (*domain.Wallet, error) { return nil, fmt.Errorf("error") }
			}

			if tt.name == "sad case - failed to update balance" {
				updateMockRepo.MockUpdateBalance = func(ctx context.Context, wallet *domain.Wallet, balance decimal.Decimal) (*domain.Wallet, error) {
					return nil, fmt.Errorf("error")
				}
			}

			wallet, err := w.DebitWallet(tt.args.ctx, tt.args.walletID, tt.args.debitAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletUsecases.DebitWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if wallet == nil {
					t.Fatalf("expected a wallet")
				}

				if wallet.Balance.String() != decimal.NewFromFloat(250).String() {
					t.Fatalf("wallet balance was not updated")
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
	amount := decimal.NewFromFloat(50)
	walletID := 2

	wal, err := w.WalletBalance(ctx, walletID)
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
				expectedBalance := wal.Balance.Add(amount)
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
