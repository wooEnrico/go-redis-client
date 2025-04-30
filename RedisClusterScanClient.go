package go_redis_scan_client

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClusterScanClient struct {
	Rdb    *redis.ClusterClient
	Option *RedisScanOptions
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
		errScan := r.BlockScan(ctx, r.Option.Match, r.Option.Count, keyChan)
		if errScan != nil {
			errChan <- errScan
		}
	}()
	return keyChan, errChan
}

func (r *RedisClusterScanClient) BlockScan(ctx context.Context, match string, count int64, channel chan string) error {
	return r.Rdb.ForEachMaster(ctx, func(context context.Context, client *redis.Client) error {
		redisScanClient := NewRedisScanClient(client, r.Option)
		return redisScanClient.BlockScan(context, match, count, channel)
	})
}
