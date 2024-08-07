package main

import (
	"fmt"
	"time"
)

// var quit = make(chan bool)

func fib(ch chan<- int, quit chan bool) {
	x, y := 1, 1

	for {
		select {
		case ch <- x:
			x, y = y, x+y
			fmt.Printf("current x:%d\n", x)
		case <-quit:
			fmt.Println("Done calculating Fibonacci!")
			close(quit)
			close(c)
			return
		}
	}
}

func main() {
	start := time.Now()

	command := ""
	ch := make(chan int)
	// ch := make(chan int, 5)
	quit := make(chan bool)

	go fib(ch, quit)

	for {
		fmt.Scanf("%s", &command)
		if command == "quit" {
			quit <- true
			break
		} else {
			fmt.Println(<-ch)
			fmt.Printf("data size:%d\n", len(ch))
		}
	}

	time.Sleep(1 * time.Second)

	elapsed := time.Since(start)
	fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}
