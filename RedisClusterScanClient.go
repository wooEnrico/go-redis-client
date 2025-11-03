package go_redis_scan_client

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClusterScanClient struct {
	Rdb    *redis.ClusterClient
	Option *RedisScanOptions
}

func (r *RedisClusterScanClient) HScan(ctx context.Context, key string) (<-chan HashTuple, <-chan error) {
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

func (r *RedisClusterScanClient) BlockHScan(ctx context.Context, key string, channel chan HashTuple) error {
	forKey, err := r.Rdb.MasterForKey(ctx, key)
	if err != nil {
		return err
	}
	redisScanClient := NewRedisScanClient(forKey, r.Option)
	return redisScanClient.BlockHScan(ctx, key, channel)
}

func (r *RedisClusterScanClient) SScan(ctx context.Context, key string) (<-chan string, <-chan error) {
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

func (r *RedisClusterScanClient) BlockSScan(ctx context.Context, key string, channel chan string) error {
	forKey, err := r.Rdb.MasterForKey(ctx, key)
	if err != nil {
		return err
	}
	redisScanClient := NewRedisScanClient(forKey, r.Option)
	return redisScanClient.BlockSScan(ctx, key, channel)
}

func (r *RedisClusterScanClient) ZScan(ctx context.Context, key string) (<-chan ZSetTuple, <-chan error) {
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

func (r *RedisClusterScanClient) BlockZScan(ctx context.Context, key string, channel chan ZSetTuple) error {
	forKey, err := r.Rdb.MasterForKey(ctx, key)
	if err != nil {
		return err
	}
	redisScanClient := NewRedisScanClient(forKey, r.Option)
	return redisScanClient.BlockZScan(ctx, key, channel)
}

func NewRedisClusterScanClient(rdb *redis.ClusterClient, option *RedisScanOptions) RedisScanClientInterface {
	return &RedisClusterScanClient{
		Rdb:    rdb,
		Option: option,
	}
}

func (r *RedisClusterScanClient) Scan(ctx context.Context) (<-chan string, <-chan error) {
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

func (r *RedisClusterScanClient) BlockScan(ctx context.Context, channel chan string) error {
	return r.Rdb.ForEachMaster(ctx, func(context context.Context, client *redis.Client) error {
		redisScanClient := NewRedisScanClient(client, r.Option)
		return redisScanClient.BlockScan(context, channel)
	})
}
