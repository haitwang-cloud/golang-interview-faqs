利用令牌桶算法实现的限流器，每秒产生固定数量的令牌，当令牌数量达到桶的容量时，多余的令牌会被丢弃。
具体的实现用golang的struct和channel来实现，代码如下：

``` go
type RateLimiter struct {
	rate         int64         // 限流速率（每秒）
	capacity     int64         // 桶的容量
	tokens       chan struct{} // 令牌通道
	tokenCounter int64         // 当前令牌数量
}
```

``` go
func (limiter *RateLimiter) tokenGenerator() {
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
```

Allow则是从令牌通道中获取令牌，如果通道为空则返回false，否则返回true。

``` go
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
```

Main函数中的测试

``` go
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
```

代码的运行结果如下：

``` shell
Request 1 allowed
Request 2 allowed
Request 3 allowed
Request 4 allowed
Request 5 allowed
Request 6 denied
Request 7 denied
Request 8 denied
Request 9 denied
Request 10 denied

```

TestRateLimiter函数的运行结果如下：

``` shell
    rateLimit_test.go:34: Request 2 allowed
    rateLimit_test.go:36: Request 0 denied
    rateLimit_test.go:34: Request 3 allowed
    rateLimit_test.go:36: Request 9 denied
    rateLimit_test.go:34: Request 8 allowed
    rateLimit_test.go:34: Request 4 allowed
    rateLimit_test.go:36: Request 6 denied
    rateLimit_test.go:36: Request 1 denied
    rateLimit_test.go:34: Request 7 allowed
    rateLimit_test.go:36: Request 5 denied
```