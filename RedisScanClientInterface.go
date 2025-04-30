package go_redis_scan_client

import "context"

type RedisScanClientInterface interface {
	Scan(ctx context.Context) (<-chan string, <-chan error)
	BlockScan(ctx context.Context, match string, count int64, channel chan string) error
}

type RedisScanOptions struct {
	ChannelSize int
	Count       int64
	Match       string
}
