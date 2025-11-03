package go_redis_scan_client

import (
	"context"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
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
		KeyType:     "",
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
		KeyType:     "",
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

func TestClusterHScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:7001", "localhost:7002", "localhost:7003"},
	})

	rdb.HSet(context.Background(), "testHash", "field1", "value1")

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	fieldChan, errChan := scan.HScan(context.Background(), "testHash")

	for field := range fieldChan {
		log.Println(field)
	}

	if err := <-errChan; err != nil {
		log.Printf("Error during HScan: %v", err)
		t.Fatal(err)
	}

	rdb.Del(context.Background(), "testHash")

	if err := <-errChan; err != nil {
		log.Printf("Error during HScan: %v", err)
		t.Fatal(err)
	}
}

func TestHScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})

	rdb.HSet(context.Background(), "testHash", "field1", "value1")

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	fieldChan, errChan := scan.HScan(context.Background(), "testHash")

	for field := range fieldChan {
		log.Println(field)
	}

	rdb.Del(context.Background(), "testHash")

	if err := <-errChan; err != nil {
		log.Printf("Error during HScan: %v", err)
		t.Fatal(err)
	}
}

func TestClusterSScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:7001", "localhost:7002", "localhost:7003"},
	})

	rdb.SAdd(context.Background(), "testSet", "member1", "member2", "member3")

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	memberChan, errChan := scan.SScan(context.Background(), "testSet")

	for member := range memberChan {
		log.Println(member)
	}

	if err := <-errChan; err != nil {
		log.Printf("Error during SScan: %v", err)
		t.Fatal(err)
	}

	rdb.Del(context.Background(), "testSet")

	if err := <-errChan; err != nil {
		log.Printf("Error during SScan: %v", err)
		t.Fatal(err)
	}
}

func TestSScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})
	rdb.SAdd(context.Background(), "testSet", "member1", "member2", "member3")

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	memberChan, errChan := scan.SScan(context.Background(), "testSet")

	for member := range memberChan {
		log.Println(member)
	}

	rdb.Del(context.Background(), "testSet")

	if err := <-errChan; err != nil {
		log.Printf("Error during SScan: %v", err)
		t.Fatal(err)
	}
}

func TestClusterZScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:7001", "localhost:7002", "localhost:7003"},
	})

	rdb.ZAdd(context.Background(), "testZSet", redis.Z{Score: 1, Member: "member1"})

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	memberChan, errChan := scan.ZScan(context.Background(), "testZSet")

	for member := range memberChan {
		log.Println(member)
	}

	if err := <-errChan; err != nil {
		log.Printf("Error during ZScan: %v", err)
		t.Fatal(err)
	}

	rdb.Del(context.Background(), "testZSet")

	if err := <-errChan; err != nil {
		log.Printf("Error during ZScan: %v", err)
		t.Fatal(err)
	}
}

func TestZScan(t *testing.T) {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})
	rdb.ZAdd(context.Background(), "testZSet", redis.Z{Score: 1, Member: "member1"})

	scan := NewRedisUniversalScanClient(rdb, &RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	memberChan, errChan := scan.ZScan(context.Background(), "testZSet")

	for member := range memberChan {
		log.Println(member)
	}

	rdb.Del(context.Background(), "testZSet")

	if err := <-errChan; err != nil {
		log.Printf("Error during ZScan: %v", err)
		t.Fatal(err)
	}
}
