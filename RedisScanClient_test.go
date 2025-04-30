package go_redis_scan_client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
)

// go test -v

func TestClusterScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:7001", "localhost:7002", "localhost:7003"},
	})

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	keyChan, errChan := scan.Scan(context.Background())

	for key := range keyChan {
		log.Println(key)
	}

	if err := <-errChan; err != nil {
		log.Printf("Error during scan: %v", err)
		t.Fatal(err)
	}
}

func TestScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	keyChan, errChan := scan.Scan(context.Background())

	for key := range keyChan {
		log.Println(key)
	}

	if err := <-errChan; err != nil {
		log.Printf("Error during scan: %v", err)
		t.Fatal(err)
	}
}
