package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// WalletDb sets up wallet's API server database layer
// with all the necessary dependencies
type WalletDb struct {
	Db  *gorm.DB
	Rdb *redis.Client
}

// NewWalletDb initializes a new wallet server database instance
// that meets all the preconsitions checks
func NewWalletDb(gorm *gorm.DB, rdb *redis.Client) *WalletDb {
	db := WalletDb{
		Db:  gorm,
		Rdb: rdb,
	}
	db.checkPreconditions()

	return &db
}

func (db *WalletDb) checkPreconditions() {
	if db.Db == nil {
		log.Panicf("error initializing database, ORM has not been initialized")
	}
}

// ConnectToDatabase opens a connection to the database
func ConnectToDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	tables := []interface{}{
		&domain.Wallet{},
	}
	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return fmt.Errorf("failed to automigrate: %v", err)
		}
	}

	return nil
}

// GetBalance retrieves a wallet balance for the supplied wallet ID
func (db *WalletDb) GetBalance(
	ctx context.Context,
	walletID int,
) (*domain.Wallet, error) {
	var wallet domain.Wallet

	id := fmt.Sprint(wallet.ID)
	result, err := db.Rdb.Get(id).Result()
	switch err {
	case redis.Nil:
		log.Printf("%s does not exist", id)

	case nil:
		if err := json.Unmarshal([]byte(result), &wallet); err != nil {
			return nil, err
		}
		return &wallet, nil

	default:
		return nil, err
	}

	if err := db.Db.First(&wallet, walletID).Error; err != nil {
		return nil, err
	}

	return &wallet, nil
}

// UpdateBalance updates (credits/debits) a wallet's balance
func (db *WalletDb) UpdateBalance(
	ctx context.Context,
	wallet *domain.Wallet,
	balance decimal.Decimal,
) (*domain.Wallet, error) {
	if wallet == nil {
		return nil, fmt.Errorf("no wallet has been passed")
	}

	if err := db.Db.Model(&wallet).
		Where("id = ?", wallet.ID).
		Updates(domain.Wallet{Balance: balance}).
		Error; err != nil {
		return nil, err
	}

	walletJson, err := json.Marshal(wallet)
	if err != nil {
		return nil, err
	}
	if err := db.Rdb.Set(fmt.Sprint(wallet.ID), walletJson, 0).Err(); err != nil {
		return nil, fmt.Errorf("error setting balance to redis: %v", err)
	}

	return wallet, nil
}
