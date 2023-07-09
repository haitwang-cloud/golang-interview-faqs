package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	rate         int64         // 限流速率（每秒）
	capacity     int64         // 桶的容量
	tokens       chan struct{} // 令牌通道
	tokenCounter int64         // 当前令牌数量
}

// NewRateLimiter 创建一个限流器
func NewRateLimiter(rate, capacity int64) *RateLimiter {
	if rate <= 0 || capacity <= 0 {
		panic("rate and capacity must be greater than 0")
	}

	limiter := &RateLimiter{
		rate:         rate,
		capacity:     capacity,
		tokens:       make(chan struct{}, capacity), // 使用有缓冲的通道
		tokenCounter: capacity,
	}

	// 初始化令牌通道
	for i := int64(0); i < capacity; i++ {
		limiter.tokens <- struct{}{}
	}

	// 启动令牌生成器
	go limiter.tokenGenerator()

	return limiter
}

// tokenGenerator 令牌生成器
func (limiter *RateLimiter) tokenGenerator() {
	// 每秒生成一个令牌
	ticker := time.NewTicker(time.Second / time.Duration(limiter.rate))
	defer ticker.Stop()

	for range ticker.C {
		// 检查当前令牌数量，如果少于桶容量，则生成新的令牌
		tokensToAdd := limiter.capacity - atomic.LoadInt64(&limiter.tokenCounter)
		for i := int64(0); i < tokensToAdd; i++ {
			// 使用非阻塞方式写入令牌通道，如果通道已满则停止生成新的令牌
			select {
			case limiter.tokens <- struct{}{}:
				atomic.AddInt64(&limiter.tokenCounter, 1)
			default:
				break
			}
		}
	}
}

// Allow 检查是否允许通过请求
func (limiter *RateLimiter) Allow() bool {
	select {
	case <-limiter.tokens:
		atomic.AddInt64(&limiter.tokenCounter, -1)
		return true // 有令牌，允许通过请求
	default:
		return false // 没有令牌，拒绝请求
	}
}

func main() {
	// 创建一个限流器，每秒生成1个令牌，桶容量为5
	limiter := NewRateLimiter(1, 5)

	// 模拟一些请求
	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "denied")
		}
	}
}
