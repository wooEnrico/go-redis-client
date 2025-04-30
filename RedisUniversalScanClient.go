package go_redis_scan_client

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type RedisUniversalScanClient struct {
	Rdb    redis.UniversalClient
	Option *RedisScanOptions
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
		errScan := r.BlockScan(ctx, r.Option.Match, r.Option.Count, keyChan)
		if errScan != nil {
			errChan <- errScan
		}
	}()
	return keyChan, errChan
}

func (r *RedisUniversalScanClient) BlockScan(ctx context.Context, match string, count int64, channel chan string) error {
	if single, ok := (r.Rdb).(*redis.Client); ok {
		return NewRedisScanClient(single, r.Option).BlockScan(ctx, match, count, channel)
	}
	if cluster, ok := (r.Rdb).(*redis.ClusterClient); ok {
		return NewRedisClusterScanClient(cluster, r.Option).BlockScan(ctx, match, count, channel)
	}
	return errors.New("unsupported Redis client type")
}
