package rdb

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// TestRedisClient 测试 RedisClient 的功能
func TestRedisClient(t *testing.T) {
	// 初始化 Redis 客户端
	redisClient, err := NewRedisClient("localhost:6379", "", 0, 10)
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}
	defer redisClient.Close()

	// 测试 Set 方法
	key := "test_key"
	value := "test_value"
	err = redisClient.Set(key, value)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}
	fmt.Printf("Set key %s successfully\n", key)

	// 测试 Get 方法
	result, err := redisClient.Get(key)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if result != value {
		t.Fatalf("Expected value %s, but got %s", value, result)
	}
	fmt.Printf("Get key %s successfully: %s\n", key, result)

	// 测试 Delete 方法
	err = redisClient.Delete(key)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}
	fmt.Printf("Deleted key %s successfully\n", key)

	// 验证键是否已删除
	result, err = redisClient.Get(key)
	if err != redis.Nil {
		t.Fatalf("Expected redis.Nil error, but got: %v", err)
	}
	fmt.Printf("Key %s is deleted as expected\n", key)
}

// TestConcurrency 测试并发情况下线程池是否生效
func TestConcurrency(t *testing.T) {
	// 初始化 Redis 客户端
	redisClient, err := NewRedisClient("localhost:6379", "", 0, 10)
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}
	defer redisClient.Close()

	// 定义任务数量
	taskCount := 10
	var wg sync.WaitGroup
	wg.Add(taskCount)

	// 记录开始时间
	startTime := time.Now()

	// 提交多个任务到线程池
	for i := 0; i < taskCount; i++ {
		key := fmt.Sprintf("concurrent_key%d", i)
		value := fmt.Sprintf("concurrent_value%d", i)

		// 提交任务到线程池
		err := redisClient.Execute(func() {
			defer wg.Done()
			err := redisClient.Set(key, value)
			if err != nil {
				t.Errorf("Failed to set key %s: %v", key, err)
			} else {
				fmt.Printf("Set key %s successfully\n", key)
			}
		})
		if err != nil {
			t.Errorf("Failed to submit task for key %s: %v", key, err)
		}
	}

	// 等待所有任务完成
	wg.Wait()

	// 计算总耗时
	elapsedTime := time.Since(startTime)
	fmt.Printf("All %d tasks completed in %v\n", taskCount, elapsedTime)

	// 验证所有键是否设置成功
	for i := 0; i < taskCount; i++ {
		key := fmt.Sprintf("concurrent_key%d", i)
		expectedValue := fmt.Sprintf("concurrent_value%d", i)

		result, err := redisClient.Get(key)
		if err != nil {
			t.Errorf("Failed to get key %s: %v", key, err)
		} else if result != expectedValue {
			t.Errorf("Expected value %s for key %s, but got %s", expectedValue, key, result)
		} else {
			fmt.Printf("Get key %s successfully: %s\n", key, result)
		}
	}
}
