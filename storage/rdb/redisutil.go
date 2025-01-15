package rdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
)

// RedisClient Redis客户端和线程池
type RedisClient struct {
	client *redis.Client      // Redis客户端
	pool   *ants.PoolWithFunc // 线程池
	ctx    context.Context    // 上下文

}

// NewRedisClient 创建RedisClient示例
func NewRedisClient(addr, password string, db int, poolSize int) (*RedisClient, error) {
	// 初始化Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,        // 设置线程池大小
		DialTimeout:  1 * time.Second, // 连接超时时间
		ReadTimeout:  1 * time.Second, // 读取超时时间
		WriteTimeout: 1 * time.Second, // 写入超时时间
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	// 初始化线程池
	pool, err := ants.NewPoolWithFunc(10, func(i interface{}) {
		task := i.(func())
		task()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create thread pool: %v", err)
	}

	return &RedisClient{
		client: client,
		pool:   pool,
		ctx:    ctx,
	}, nil
}

// Close 关闭Redis客户端和线程池
func (rc *RedisClient) Close() {
	rc.client.Close()
	rc.pool.Release()
}

// Execute 将任务提交到线程池中执行
func (rc *RedisClient) Execute(task func()) error {
	return rc.pool.Invoke(task)
}

// Set 设置缓存
func (rc *RedisClient) Set(key string, value interface{}) error {
	var wg sync.WaitGroup
	wg.Add(1)

	err := rc.Execute(func() {
		defer wg.Done()
		if err := rc.client.Set(rc.ctx, key, value, 60).Err(); err != nil {
			fmt.Printf("failed to set key %s: %v\n", key, err)
		}
	})

	wg.Wait()
	return err
}

// Get 获取缓存
func (rc *RedisClient) Get(key string) (string, error) {
	var result string
	var err error
	var wg sync.WaitGroup
	wg.Add(1)

	taskErr := rc.Execute(func() {
		defer wg.Done()
		result, err = rc.client.Get(rc.ctx, key).Result()
		if err != nil {
			fmt.Printf("failed to get key %s: %v\n", key, err)
		}
	})

	wg.Wait()
	if taskErr != nil {
		return "", taskErr
	}
	return result, err
}

// Delete 删除键
func (rc *RedisClient) Delete(key string) error {
	var wg sync.WaitGroup
	wg.Add(1)

	err := rc.Execute(func() {
		defer wg.Done()
		if err := rc.client.Del(rc.ctx, key).Err(); err != nil {
			fmt.Printf("failed to delete key %s: %v\n", key, err)
		}
	})

	wg.Wait()
	return err
}
