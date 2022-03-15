package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/ageeknamedslickback/wallet-API/wallet/dto"
	"github.com/go-redis/redis"
)

type WalletCache interface {
	CacheBalance(
		ctx context.Context,
		wallet *domain.Wallet,
	) (*domain.Wallet, error)
	GetCachedBalance(
		ctx context.Context,
		walletID int,
	) (*domain.Wallet, error)
}

// ServiceCache sets up wallet's API server cache layer
// with all the necessary dependencies
type ServiceCache struct {
	Rdb *redis.Client
}

// NewCacheService initalizes a new cache service
func NewCacheService(client *redis.Client) *ServiceCache {
	c := &ServiceCache{
		Rdb: client,
	}
	c.checkPreconditions()
	return c
}

func (c *ServiceCache) checkPreconditions() {
	if c.Rdb == nil {
		log.Panicf("cache service has not initalized redis client")
	}
}

// CacheBalance caches a wallet to easily retrieve its balance
func (c *ServiceCache) CacheBalance(
	ctx context.Context,
	wallet *domain.Wallet,
) (*domain.Wallet, error) {
	if wallet == nil {
		return nil, dto.Wrap(fmt.Errorf("no wallet has been passed"), "CacheBalance")
	}

	bs, err := json.Marshal(wallet)
	if err != nil {
		return nil, dto.Wrap(
			fmt.Errorf("failed to marshal wallet balance with err %v", err),
			"CacheBalance",
		)
	}
	if err := c.Rdb.Set(fmt.Sprint(wallet.ID), bs, 0).Err(); err != nil {
		return nil, dto.Wrap(
			fmt.Errorf("failed to cache wallet balance with err %v", err),
			"CacheBalance",
		)
	}

	return wallet, nil
}

// GetCachedBalance retrieves wallet balance from the cache
func (c *ServiceCache) GetCachedBalance(
	ctx context.Context,
	walletID int,
) (*domain.Wallet, error) {
	var wallet domain.Wallet
	result, err := c.Rdb.Get(fmt.Sprint(walletID)).Result()
	switch err {
	case redis.Nil:
		return nil, nil

	case nil:
		if err := json.Unmarshal([]byte(result), &wallet); err != nil {
			return nil, dto.Wrap(
				fmt.Errorf(
					"failed to unmarshal cached balance with err %v",
					err,
				),
				"GetCachedBalance",
			)
		}
		return &wallet, nil

	default:
		return nil, dto.Wrap(
			fmt.Errorf("failed to get cached balance with err %v", err),
			"GetCachedBalance",
		)
	}
}
