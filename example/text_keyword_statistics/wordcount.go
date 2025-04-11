package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// 配置参数
const (
	WorkerCount = 100              // 并发处理的协程数
	BufferSize  = 10 * 1024 * 1024 // 每次读取的块大小（MB）
)

// processChunk 处理文件分片
func processChunk(chunk []byte, word []byte, wg *sync.WaitGroup, resultChan chan int) {
	defer wg.Done()
	count := bytes.Count(chunk, word) // 子串统计
	resultChan <- count
}

func main() {
	/*
		实现思路：
		1、使用流式读取文件，避免一次性加载到内存。
		2、将文本按行或块读取，送入多个 goroutine 并发处理。
		3、统计目标词的出现次数，结果通过 channel 汇总。
	*/
	startTime := time.Now()

	if len(os.Args) < 3 {
		fmt.Println("用法：go run main.go <文件路径> <目标词>")
		return
	}

	filename := os.Args[1]
	targetWord := []byte(os.Args[2])
	overlap := len(targetWord) - 1 // 为了防止跨块漏匹配

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("打开文件失败：%v\n", err)
		return
	}
	defer file.Close()

	resultChan := make(chan int, WorkerCount)
	totalCount := 0

	// 协程池
	pool, err := ants.NewPool(10, ants.WithNonblocking(false), ants.WithExpiryDuration(10*time.Second))
	if err != nil {
		fmt.Printf("协程池初始化失败：%v\n", err)
		return
	}
	defer pool.Release()

	var wg sync.WaitGroup
	var prevTail []byte
	reader := bufio.NewReader(file)

	for {
		buffer := make([]byte, BufferSize)
		n, err := reader.Read(buffer)
		if n > 0 {
			// 拼接上一次块的尾部，防止目标词被拆开
			chunk := append([]byte{}, prevTail...)
			chunk = append(chunk, buffer[:n]...)

			// 更新尾部，用于下一次拼接
			if n >= overlap {
				prevTail = buffer[n-overlap : n]
			} else {
				prevTail = buffer[:n]
			}

			wg.Add(1)

			// 在提交任务之前，打印当前池子的状态
			//fmt.Printf("当前池子状态 - 容量: %d, 正在运行: %d, 可用任务: %d\n", pool.Cap(), pool.Running(), pool.Free())

			// 提交到协程池
			tmpChunk := chunk // 避免闭包中 chunk 被复用
			err := pool.Submit(func() {
				processChunk(tmpChunk, targetWord, &wg, resultChan)
			})
			if err != nil {
				fmt.Printf("提交任务失败: %v\n", err)
				wg.Done() // 保证不会永久阻塞
			}
		}

		if err != nil {
			if err == io.EOF {
				break // 文件读取完毕
			}
			fmt.Printf("读取文件失败: %v\n", err)
			break
		}
	}

	// 关闭协程
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for count := range resultChan {
		totalCount += count
	}

	fmt.Printf("目标词 `%s` 在文件中出现了 `%d` 次\n", targetWord, totalCount)
	fmt.Printf("程序执行耗时: %s\n", time.Since(startTime))
}
