package go_redis_scan_client

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisScanClient struct {
	Rdb    *redis.Client
	Option *RedisScanOptions
}

func (r *RedisScanClient) HScan(ctx context.Context, key string) (<-chan HashTuple, <-chan error) {
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

func (r *RedisScanClient) BlockHScan(ctx context.Context, key string, channel chan HashTuple) error {
	keyType, err := r.Rdb.Type(ctx, key).Result()
	if err != nil || keyType != "hash" {
		return nil
	}

	cursor := uint64(0)
	for {
		fields, nextCursor, errScan := r.Rdb.HScan(ctx, key, cursor, r.Option.Match, r.Option.Count).Result()
		if errScan != nil {
			return errScan
		}

		for i := 0; i < len(fields); i += 2 {
			field := fields[i]
			value := fields[i+1]
			fieldTuple := HashTuple{
				Field: field,
				Value: value,
			}
			select {
			case channel <- fieldTuple:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}

func (r *RedisScanClient) SScan(ctx context.Context, key string) (<-chan string, <-chan error) {
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

func (r *RedisScanClient) BlockSScan(ctx context.Context, key string, channel chan string) error {
	keyType, err := r.Rdb.Type(ctx, key).Result()
	if err != nil || keyType != "set" {
		return nil
	}
	cursor := uint64(0)
	for {
		members, nextCursor, errScan := r.Rdb.SScan(ctx, key, cursor, r.Option.Match, r.Option.Count).Result()
		if errScan != nil {
			return errScan
		}
		for _, member := range members {
			select {
			case channel <- member:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}

func (r *RedisScanClient) ZScan(ctx context.Context, key string) (<-chan ZSetTuple, <-chan error) {
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

func (r *RedisScanClient) BlockZScan(ctx context.Context, key string, channel chan ZSetTuple) error {
	keyType, err := r.Rdb.Type(ctx, key).Result()
	if err != nil || keyType != "zset" {
		return nil
	}
	cursor := uint64(0)
	for {
		members, nextCursor, errScan := r.Rdb.ZScan(ctx, key, cursor, r.Option.Match, r.Option.Count).Result()
		if errScan != nil {
			return errScan
		}

		for i := 0; i < len(members); i += 2 {
			member := members[i]
			scoreStr := members[i+1]
			score, err := strconv.ParseFloat(scoreStr, 64)
			if err != nil {
				return err
			}
			memberTuple := ZSetTuple{
				Member: member,
				Score:  score,
			}
			select {
			case channel <- memberTuple:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
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
		errScan := r.BlockScan(ctx, keyChan)
		if errScan != nil {
			errChan <- errScan
		}
	}()
	return keyChan, errChan
}

func (r *RedisScanClient) BlockScan(ctx context.Context, channel chan string) error {
	var cursor uint64 = 0
	for {
		keys, nextCursor, errScan := r.Rdb.ScanType(ctx, cursor, r.Option.Match, r.Option.Count, r.Option.KeyType).Result()
		if errScan != nil {
			return errScan
		}

		for _, key := range keys {
			select {
			case channel <- key:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}
