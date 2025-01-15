package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/wenhao26/toolbox/storage/rdb"
)

func main() {
	// 初始化 Redis 客户端
	redisClient, err := rdb.NewRedisClient("localhost:6379", "", 9, 10)
	if err != nil {
		fmt.Println("Failed to create Redis client:", err)
		return
	}
	defer redisClient.Close()

	// 定义任务数量
	taskCount := 100
	var wg sync.WaitGroup
	wg.Add(taskCount)

	// 记录开始时间
	startTime := time.Now()

	// 提交多个任务到线程池
	for i := 0; i < taskCount; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)

		go func(k, v string, taskID int) {
			defer wg.Done()

			// 记录任务开始时间
			taskStartTime := time.Now()
			fmt.Printf("Task %d started at %v\n", taskID, taskStartTime)

			// 模拟任务执行时间
			time.Sleep(1 * time.Second)

			// 执行 Redis 操作
			err := redisClient.Set(k, v)
			if err != nil {
				fmt.Printf("Task %d failed to set key %s: %v\n", taskID, k, err)
			} else {
				fmt.Printf("Task %d set key %s successfully\n", taskID, k)
			}

			// 记录任务结束时间
			taskEndTime := time.Now()
			fmt.Printf("Task %d finished at %v (elapsed: %v)\n", taskID, taskEndTime, taskEndTime.Sub(taskStartTime))
		}(key, value, i)
	}

	// 等待所有任务完成
	wg.Wait()

	// 计算总耗时
	elapsedTime := time.Since(startTime)
	fmt.Printf("All tasks completed in %v\n", elapsedTime)
}
