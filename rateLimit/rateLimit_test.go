package main

import (
	"sync"
	"testing"
	"time"
)

const (
	defaultReqNum   = 10
	defaultCapacity = 5
	defaultRate     = 1

	defaultWaitTime = 100 * time.Millisecond
)

func TestRateLimiter(t *testing.T) {
	// 建一个限流器，每秒生成1个令牌，桶容量为5
	limiter := NewRateLimiter(defaultRate, defaultCapacity)

	// 使用等待组来同步并发测试
	var wg sync.WaitGroup
	wg.Add(defaultReqNum)

	// 并发模拟10个请求
	for i := 0; i < defaultReqNum; i++ {
		go func(index int) {
			defer wg.Done()

			// 等待100毫秒，以模拟请求之间的间隔
			time.Sleep(defaultWaitTime)

			if limiter.Allow() {
				t.Logf("Request %d allowed", index)
			} else {
				t.Logf("Request %d denied", index)
			}
		}(i)
	}

	// 等待所有请求完成
	wg.Wait()
}
