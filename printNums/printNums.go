package main

import (
	"fmt"
	"sync"
)

/*
   使用两个channel分别传递奇数和偶数
   printNumbers goroutine在两个channel间轮询,接收数后打印
   main函数先在oddCh发送一个数来启动打印
   最后使用waitGroup来等待打印结束
*/
func printNumbers(oddCh, evenCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	i := 1
	for i <= 10 {
		select {
		case num := <-oddCh:
			fmt.Println("Odd:", num)
			i += 1
			evenCh <- i
		case num := <-evenCh:
			fmt.Println("Even:", num)
			i += 1
			oddCh <- i
		}
	}
}

func main() {
	oddCh := make(chan int, 1)  // 带有缓冲的通道
	evenCh := make(chan int, 1) // 带有缓冲的通道
	var wg sync.WaitGroup
	wg.Add(1) // 只有一个goroutine需要等待

	go func() {
		printNumbers(oddCh, evenCh, &wg)
	}()

	oddCh <- 1 // 启动奇数打印
	wg.Wait()
}
