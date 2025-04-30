# Redis Scan Client

一个用于扫描 Redis 键的 Go 客户端库，支持单机和集群模式。

## 功能特点

- 支持单机和集群模式的 Redis 扫描。
- 异步扫描，使用通道返回结果。
- 支持自定义匹配模式和扫描数量。

## API 文档

### RedisScanOptions

| 参数名      | 类型   | 描述                     |
|-------------|--------|--------------------------|
| ChannelSize | int    | 通道缓冲区大小           |
| Count       | int64  | 每次扫描的键数量         |
| Match       | string | 键的匹配模式，支持通配符 |

### RedisScanClient

- `NewRedisUniversalScanClient(rdb redis.UniversalClient, option *RedisScanOptions)`：创建新的扫描客户端。
- `Scan(ctx context.Context) (<-chan string, <-chan error)`：开始扫描，返回键通道和异常通道。


## 安装

```bash
go get github.com/wooenrico/go-redis-client
```

## 使用示例

### 单机模式

```go
package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	goRedisScan "github.com/wooenrico/go-redis-client"
	"log"
)

func main() {
	// 创建Redis客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})
	// 创建扫描客户端
	scan := goRedisScan.NewRedisUniversalScanClient(rdb, &goRedisScan.RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	keyChan, errChan := scan.Scan(context.Background())

	for key := range keyChan {
		log.Println(key)
	}

	// 判断是否因为异常扫描结束
	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}

```

### 集群模式

```go
package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	goRedisScan "github.com/wooenrico/go-redis-client"
	"log"
)

func main() {
	// 创建Redis 集群客户端
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:7001", "localhost:7002", "localhost:7003"},
	})
	// 创建扫描客户端
	scan := goRedisScan.NewRedisUniversalScanClient(rdb, &goRedisScan.RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*",
	})
	// 开始扫描
	keyChan, errChan := scan.Scan(context.Background())

	for key := range keyChan {
		log.Println(key)
	}

	// 判断是否因为异常扫描结束
	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}

```

## 贡献

欢迎提交 Issue 和 Pull Request！
