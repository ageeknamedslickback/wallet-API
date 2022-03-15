package cache_test

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/ageeknamedslickback/wallet-API/wallet/domain"
	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/services/cache"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
)

var ctx = context.Background()

func initalizeRedisService() *cache.ServiceCache {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return cache.NewCacheService(rdb)
}
func TestServiceCache_CacheBalance(t *testing.T) {
	type args struct {
		ctx    context.Context
		wallet *domain.Wallet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: ctx,
				wallet: &domain.Wallet{
					ID:      10,
					Balance: decimal.NewFromFloat(10.45),
				},
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
			c := initalizeRedisService()
			wallet, err := c.CacheBalance(tt.args.ctx, tt.args.wallet)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceCache.CacheBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && wallet == nil {
				t.Fatalf("expected a cached wallet balance")
			}

			if tt.wantErr && wallet != nil {
				t.Fatalf("expected no cached wallet balance")
			}
		})
	}
}

func TestServiceCache_GetCachedBalance(t *testing.T) {
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
				walletID: 10,
			},
			wantErr: false,
		},
		{
			name: "sad case; redis nil",
			args: args{
				ctx:      ctx,
				walletID: 100000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := initalizeRedisService()

			wallet, err := c.GetCachedBalance(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceCache.GetCachedBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && wallet != nil && tt.name == "sad case; redis nil" {
				t.Fatalf("expected no cached wallet")
			}

			if !tt.wantErr && wallet == nil && tt.name == "sad case; failed to get cached balance" {
				t.Fatalf("expected a cached wallet balance")
			}

			if tt.wantErr && wallet != nil {
				t.Fatalf("expected no cached wallet balance")
			}

			os.Setenv("REDIS_ADDR", os.Getenv("REDIS_ADDR"))
		})
	}
}
