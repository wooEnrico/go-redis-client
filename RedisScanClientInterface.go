package go_redis_scan_client

import "context"

type RedisScanClientInterface interface {

	// Scan 扫描所有键
	Scan(ctx context.Context) (<-chan string, <-chan error)
	// BlockScan 阻塞扫描所有键
	BlockScan(ctx context.Context, channel chan string) error
	// HScan 扫描哈希表
	HScan(ctx context.Context, key string) (<-chan HashTuple, <-chan error)
	// BlockHScan 阻塞扫描哈希表
	BlockHScan(ctx context.Context, key string, channel chan HashTuple) error
	// SScan 扫描集合
	SScan(ctx context.Context, key string) (<-chan string, <-chan error)
	// BlockSScan 阻塞扫描集合
	BlockSScan(ctx context.Context, key string, channel chan string) error
	// ZScan 扫描有序集合
	ZScan(ctx context.Context, key string) (<-chan ZSetTuple, <-chan error)
	// BlockZScan 阻塞扫描有序集合
	BlockZScan(ctx context.Context, key string, channel chan ZSetTuple) error
}

// RedisScanOptions 扫描选项
type RedisScanOptions struct {
	ChannelSize int
	Count       int64
	Match       string
	KeyType     string
}

// HashTuple 哈希表扫描返回的结构体，包含字段和值
type HashTuple struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// ZSetTuple 有序集合扫描返回的结构体，包含成员和分数
type ZSetTuple struct {
	Member string  `json:"member"`
	Score  float64 `json:"score"`
}
