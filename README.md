# Redis Scan Client

一个用于扫描 Redis 键的 Go 客户端库，支持单机和集群模式。

## 功能特点

- 支持单机和集群模式的 Redis 扫描
- 异步扫描，使用通道返回结果
- 支持阻塞和非阻塞两种扫描模式
- 支持键（Key）扫描
- 支持哈希表（Hash）字段扫描，返回字段和值的键值对
- 支持集合（Set）成员扫描
- 支持有序集合（Sorted Set）成员扫描，返回成员和分数
- 支持自定义匹配模式和扫描数量
- 支持按数据类型过滤扫描键

## API 文档

### RedisScanOptions

| 参数名      | 类型   | 描述                           |
|-------------|--------|--------------------------------|
| ChannelSize | int    | 通道缓冲区大小                 |
| Count       | int64  | 每次扫描的键数量               |
| Match       | string | 键的匹配模式，支持通配符       |
| KeyType     | string | 键的类型过滤，如 "hash"、"set"、"zset" 等，为空表示不过滤 |

### 数据结构

#### HashTuple

哈希表扫描返回的结构体，包含字段和值：

```go
type HashTuple struct {
    Field string  // 哈希表字段名
    Value string  // 哈希表字段值
}
```

#### ZSetTuple

有序集合扫描返回的结构体，包含成员和分数：

```go
type ZSetTuple struct {
    Member string  // 有序集合成员
    Score  float64 // 成员分数
}
```

### RedisScanClientInterface

所有扫描客户端都实现了 `RedisScanClientInterface` 接口，提供以下方法：

#### 键扫描方法

- `Scan(ctx context.Context) (<-chan string, <-chan error)`：异步扫描所有键，返回键通道和错误通道
- `BlockScan(ctx context.Context, channel chan string) error`：阻塞扫描所有键，结果通过提供的通道返回

#### 哈希表扫描方法

- `HScan(ctx context.Context, key string) (<-chan HashTuple, <-chan error)`：异步扫描哈希表，返回 `HashTuple` 通道和错误通道
- `BlockHScan(ctx context.Context, key string, channel chan HashTuple) error`：阻塞扫描哈希表，返回 `HashTuple` 结构体（包含字段和值）

#### 集合扫描方法

- `SScan(ctx context.Context, key string) (<-chan string, <-chan error)`：异步扫描集合成员，返回成员通道和错误通道
- `BlockSScan(ctx context.Context, key string, channel chan string) error`：阻塞扫描集合成员，结果通过提供的通道返回

#### 有序集合扫描方法

- `ZScan(ctx context.Context, key string) (<-chan ZSetTuple, <-chan error)`：异步扫描有序集合，返回 `ZSetTuple` 通道和错误通道
- `BlockZScan(ctx context.Context, key string, channel chan ZSetTuple) error`：阻塞扫描有序集合，返回 `ZSetTuple` 结构体（包含成员和分数）

### 客户端创建函数

- `NewRedisScanClient(rdb *redis.Client, option *RedisScanOptions) RedisScanClientInterface`：创建单机 Redis 扫描客户端
- `NewRedisClusterScanClient(rdb *redis.ClusterClient, option *RedisScanOptions) RedisScanClientInterface`：创建集群 Redis 扫描客户端
- `NewRedisUniversalScanClient(rdb redis.UniversalClient, option *RedisScanOptions) RedisScanClientInterface`：创建通用 Redis 扫描客户端（自动识别单机或集群模式）


## 安装

```bash
go get github.com/wooenrico/go-redis-client
```

## 使用示例

### 键扫描

#### 单机模式

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
		KeyType:     "", // 为空表示扫描所有类型的键
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

#### 集群模式

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

### 哈希表扫描

```go
package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	goRedisScan "github.com/wooenrico/go-redis-client"
	"log"
)

func main() {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})

	scan := goRedisScan.NewRedisUniversalScanClient(rdb, &goRedisScan.RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*", // 匹配字段名模式
	})

	// 扫描哈希表，返回 HashTuple（包含字段和值）
	fieldChan, errChan := scan.HScan(context.Background(), "myHash")

	for tuple := range fieldChan {
		log.Printf("Field: %s, Value: %s", tuple.Field, tuple.Value)
	}

	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}
```

### 集合扫描

```go
package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	goRedisScan "github.com/wooenrico/go-redis-client"
	"log"
)

func main() {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})

	scan := goRedisScan.NewRedisUniversalScanClient(rdb, &goRedisScan.RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*", // 匹配成员名模式
	})

	// 扫描集合成员
	memberChan, errChan := scan.SScan(context.Background(), "mySet")

	for member := range memberChan {
		log.Println("Member:", member)
	}

	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}
```

### 有序集合扫描

```go
package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	goRedisScan "github.com/wooenrico/go-redis-client"
	"log"
)

func main() {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	})

	scan := goRedisScan.NewRedisUniversalScanClient(rdb, &goRedisScan.RedisScanOptions{
		ChannelSize: 100,
		Count:       100,
		Match:       "*", // 匹配成员名模式
	})

	// 扫描有序集合，返回 ZSetTuple（包含成员和分数）
	memberChan, errChan := scan.ZScan(context.Background(), "myZSet")

	for tuple := range memberChan {
		log.Printf("Member: %s, Score: %.2f", tuple.Member, tuple.Score)
	}

	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}
```


## 注意事项

### 错误处理

- 异步扫描方法（`Scan`、`HScan`、`SScan`、`ZScan`）返回两个通道：结果通道和错误通道
- 错误通道会在扫描完成或出错时发送错误信息，之后通道会被关闭
- **必须**从错误通道读取一次，以确认扫描是否正常完成
- 如果错误通道返回 `nil`，表示扫描成功完成
- 如果错误通道返回非 `nil` 错误，表示扫描过程中出现了错误

### 数据结构说明

- **HashTuple**：哈希表扫描返回的结构体，包含 `Field`（字段名）和 `Value`（字段值）两个字段，可以直接访问
- **ZSetTuple**：有序集合扫描返回的结构体，包含 `Member`（成员）和 `Score`（分数）两个字段，可以直接访问
- **SScan**：集合扫描直接返回字符串类型的成员名

### Context 取消

- 所有方法都支持通过 `context.Context` 进行取消操作
- 当 context 被取消时，扫描会立即停止并返回 `context.Canceled` 错误
- 这对于实现超时控制非常有用

### 匹配模式

- `Match` 参数支持 Redis 通配符模式：
  - `*`：匹配任意字符
  - `?`：匹配单个字符
  - `[abc]`：匹配指定字符中的一个
  - 空字符串表示不进行匹配过滤

### 类型过滤

- `KeyType` 参数用于在扫描键时按类型过滤，支持的值：
  - `"string"`：字符串类型
  - `"hash"`：哈希表类型
  - `"list"`：列表类型
  - `"set"`：集合类型
  - `"zset"`：有序集合类型
  - 空字符串表示不过滤，扫描所有类型的键

## 最佳实践

### 通道缓冲区大小

- `ChannelSize` 应根据数据量和处理速度合理设置
- 过小的缓冲区可能导致发送方阻塞
- 过大的缓冲区可能占用过多内存
- 建议根据实际情况调整，一般设置为 100-1000

### 扫描数量

- `Count` 参数控制每次扫描返回的元素数量（近似值）
- 较大的值可以减少网络往返次数，但会增加单次扫描的时间
- 较小的值可以提高响应速度，但会增加网络往返次数
- 建议设置为 100-1000，具体取决于数据大小和网络延迟

### 使用 Context 超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

keyChan, errChan := scan.Scan(ctx)

// 处理结果...
```

### 资源清理

确保在程序退出时正确关闭 Redis 客户端：

```go
defer rdb.Close()
```

## 版本要求

- Go 1.21 或更高版本
- Redis Server 2.8.0 或更高版本（支持 SCAN 命令）

## 依赖

- `github.com/redis/go-redis/v9` - Redis 客户端库

## 贡献

欢迎提交 Issue 和 Pull Request！
