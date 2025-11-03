package go_redis_scan_client

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrUnsupportedRedisType = errors.New("unsupported Redis client type")

type RedisUniversalScanClient struct {
	Rdb    redis.UniversalClient
	Option *RedisScanOptions
}

func (r *RedisUniversalScanClient) HScan(ctx context.Context, key string) (<-chan HashTuple, <-chan error) {
	keyChan := make(chan HashTuple, r.Option.ChannelSize)
	errChan := make(chan error, 1)
	go func() {
		defer close(keyChan)
		defer close(errChan)
		errScan := r.BlockHScan(ctx, key, keyChan)
		if errScan != nil {
			errChan <- errScan
		}
	}()
	return keyChan, errChan
}

func (r *RedisUniversalScanClient) BlockHScan(ctx context.Context, key string, channel chan HashTuple) error {
	if single, ok := (r.Rdb).(*redis.Client); ok {
		return NewRedisScanClient(single, r.Option).BlockHScan(ctx, key, channel)
	}
	if cluster, ok := (r.Rdb).(*redis.ClusterClient); ok {
		return NewRedisClusterScanClient(cluster, r.Option).BlockHScan(ctx, key, channel)
	}
	return ErrUnsupportedRedisType
}

func (r *RedisUniversalScanClient) SScan(ctx context.Context, key string) (<-chan string, <-chan error) {
	keyChan := make(chan string, r.Option.ChannelSize)
	errChan := make(chan error, 1)
	go func() {
		defer close(keyChan)
		defer close(errChan)
		errScan := r.BlockSScan(ctx, key, keyChan)
		if errScan != nil {
			errChan <- errScan
		}
	}()
	return keyChan, errChan
}

func (r *RedisUniversalScanClient) BlockSScan(ctx context.Context, key string, channel chan string) error {
	if single, ok := (r.Rdb).(*redis.Client); ok {
		return NewRedisScanClient(single, r.Option).BlockSScan(ctx, key, channel)
	}
	if cluster, ok := (r.Rdb).(*redis.ClusterClient); ok {
		return NewRedisClusterScanClient(cluster, r.Option).BlockSScan(ctx, key, channel)
	}
	return ErrUnsupportedRedisType
}

func (r *RedisUniversalScanClient) ZScan(ctx context.Context, key string) (<-chan ZSetTuple, <-chan error) {
	keyChan := make(chan ZSetTuple, r.Option.ChannelSize)
	errChan := make(chan error, 1)
	go func() {
		defer close(keyChan)
		defer close(errChan)
		errScan := r.BlockZScan(ctx, key, keyChan)
		if errScan != nil {
			errChan <- errScan
		}
	}()
	return keyChan, errChan
}

func (r *RedisUniversalScanClient) BlockZScan(ctx context.Context, key string, channel chan ZSetTuple) error {
	if single, ok := (r.Rdb).(*redis.Client); ok {
		return NewRedisScanClient(single, r.Option).BlockZScan(ctx, key, channel)
	}
	if cluster, ok := (r.Rdb).(*redis.ClusterClient); ok {
		return NewRedisClusterScanClient(cluster, r.Option).BlockZScan(ctx, key, channel)
	}
	return ErrUnsupportedRedisType
}

func NewRedisUniversalScanClient(rdb redis.UniversalClient, option *RedisScanOptions) RedisScanClientInterface {
	return &RedisUniversalScanClient{
		Rdb:    rdb,
		Option: option,
	}
}

func (r *RedisUniversalScanClient) Scan(ctx context.Context) (<-chan string, <-chan error) {
	keyChan := make(chan string, r.Option.ChannelSize)
	errChan := make(chan error, 1)
	go func() {
		defer close(keyChan)
		defer close(errChan)
		errScan := r.BlockScan(ctx, keyChan)
		if errScan != nil {
			errChan <- errScan
		}
	}()
	return keyChan, errChan
}

func (r *RedisUniversalScanClient) BlockScan(ctx context.Context, channel chan string) error {
	if single, ok := (r.Rdb).(*redis.Client); ok {
		return NewRedisScanClient(single, r.Option).BlockScan(ctx, channel)
	}
	if cluster, ok := (r.Rdb).(*redis.ClusterClient); ok {
		return NewRedisClusterScanClient(cluster, r.Option).BlockScan(ctx, channel)
	}
	return ErrUnsupportedRedisType
}
