package database_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/database"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/shopspring/decimal"
)

func initTestDatabase() *database.WalletDb {
	gormDb, err := database.ConnectToDatabase()
	if err != nil {
		log.Panicf("error connecting to the database: %v", err)
	}

	return database.NewWalletDb(gormDb)
}

func TestWalletDb_GetWallet(t *testing.T) {
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
			db := initTestDatabase()
			wallet, err := db.GetWallet(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"WalletDb.GetWallet() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if !tt.wantErr {
				if wallet == nil {
					t.Fatalf("expected a wallet")
				}

				if wallet.Balance.String() != decimal.NewFromInt(100).String() {
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

func TestWalletDb_UpdateBalance(t *testing.T) {
	ctx := context.Background()
	db := initTestDatabase()

	walletID := 2 // existing wallet
	wallet, err := db.GetWallet(ctx, walletID)
	if err != nil {
		t.Fatalf("expected to get wallet with id 2: %v", err)
	}

	newBalance := decimal.NewFromFloat(50.32)

	type args struct {
		ctx     context.Context
		wallet  *domain.Wallet
		balance decimal.Decimal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				wallet:  wallet,
				balance: newBalance,
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:    ctx,
				wallet: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := db.UpdateBalance(
				tt.args.ctx,
				tt.args.wallet,
				tt.args.balance,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"WalletDb.UpdateBalance() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if !tt.wantErr {
				if wallet == nil {
					t.Fatalf("expected wallet to be returned")
				}
				if wallet.Balance != newBalance {
					t.Fatalf(
						"expected a balance of 50.32 but got %s",
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

func TestConnectToDatabase(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "happy case",
			wantErr: false,
		},
		{
			name:    "sad case - non existent database",
			wantErr: true,
		},
		{
			name:    "sad case - wrong user password",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case - non existent database" {
				os.Setenv("DB_NAME", gofakeit.Name())
			}

			if tt.name == "sad case - wrong user password" {
				os.Setenv("DB_PASS", gofakeit.FarmAnimal())
			}

			db, err := database.ConnectToDatabase()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ConnectToDatabase() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr && db == nil {
				t.Fatalf("expected a *gorm.DB object")
			}
		})
	}
}
