总体思路是:
    使用两个channel分别传递奇数和偶数
    printNumbers goroutine在两个channel间轮询,接收数后打印
    main函数先在oddCh发送一个数来启动打印
    最后使用waitGroup来等待打印结束

```
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
```

相比直接使用goroutine,这种channel方式可以更好地控制和同步goroutine之间的执行顺序。

主要优点是:

    channel具备队列,可以缓冲数据
    select可以同时等待多个channel,实现多路复用
    不需要通过共享变量和lock来同步不同goroutine

代码输出结果如下:

```
Odd: 1
Even: 2
Odd: 3
Even: 4
Odd: 5
Even: 6
Odd: 7
Even: 8
Odd: 9
Even: 10
```