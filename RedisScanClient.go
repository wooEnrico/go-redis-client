package go_redis_scan_client

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisScanClient struct {
	Rdb    *redis.Client
	Option *RedisScanOptions
}

func NewRedisScanClient(rdb *redis.Client, option *RedisScanOptions) RedisScanClientInterface {
	return &RedisScanClient{
		Rdb:    rdb,
		Option: option,
	}
}

func (r *RedisScanClient) Scan(ctx context.Context) (<-chan string, <-chan error) {
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

func (r *RedisScanClient) BlockScan(ctx context.Context, match string, count int64, channel chan string) error {
	var cursor uint64 = 0
	for {
		keys, nextCursor, errScan := r.Rdb.Scan(ctx, cursor, match, count).Result()
		if errScan != nil {
			return errScan
		}

		for _, key := range keys {
			channel <- key
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}
