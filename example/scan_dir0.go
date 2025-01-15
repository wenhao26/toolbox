package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func walkDir(path string, wg *sync.WaitGroup, fileSizes chan<- int64) {
	defer wg.Done()

	for _, f := range getFiles(path) {
		if f.IsDir() {
			childDir := filepath.Join(path, f.Name())
			wg.Add(1)
			go walkDir(childDir, wg, fileSizes)
		} else {
			// fmt.Printf(">>>Name=%s,Size=%v\n", f.Name(), f.Size())
			fileSizes <- f.Size()
		}
	}
}

var handler = make(chan struct{}, 30) // 最多同时读取目录的协程数

func getFiles(path string) []fs.FileInfo {
	handler <- struct{}{}
	defer func() {
		<-handler
	}()

	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "err:%v\n", err)
		return nil
	}
	return fileInfos
}

func main() {
	sTime := time.Now()
	path := "D:\\www\\dev"

	fileSizes := make(chan int64)
	var wg sync.WaitGroup

	wg.Add(1)
	go walkDir(path, &wg, fileSizes)

	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	// Print the results.
	var nFiles, nBytes int64
	for size := range fileSizes {
		nFiles++
		nBytes += size
	}

	fmt.Printf("Files=%d Size=%.1f GB Time-Consuming=%v\n", nFiles, float64(nBytes)/1e9, time.Since(sTime))
}
