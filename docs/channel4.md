# Channel 作業解答

使用並行快速地計算費氏數值
您在改良後的程式版本中實作並行，加快了程式的執行速度，如下所示：

```Go


package main

import (
    "fmt"
    "math/rand"
    "time"
)

func fib(number float64, ch chan string) {
    x, y := 1.0, 1.0
    for i := 0; i < int(number); i++ {
        x, y = y, x+y
    }

    r := rand.Intn(3)
    time.Sleep(time.Duration(r) * time.Second)

    ch <- fmt.Sprintf("Fib(%v): %v\n", number, x)
}

func main() {
    start := time.Now()

    size := 15
    ch := make(chan string, size)

    for i := 0; i < size; i++ {
        go fib(float64(i), ch)
    }

    for i := 0; i < size; i++ {
        fmt.Printf(<-ch)
    }

    elapsed := time.Since(start)
    fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}
```
您在第二版的程式中使用了兩個無法緩衝的通道，如下所示：

```Go

package main

import (
    "fmt"
    "time"
)

var quit = make(chan bool)

func fib(c chan int) {
    x, y := 1, 1

    for {
        select {
            case c <- x:
                x, y = y, x+y
            case <-quit:
                fmt.Println("Done calculating Fibonacci!")
            return
        }
    }
}

func main() {
    start := time.Now()

    command := ""
    data := make(chan int)

    go fib(data)

    for {
        num := <-data
        fmt.Println(num)
        fmt.Scanf("%s", &command)
        if command == "quit" {
            quit <- true
            break
        }
    }

    time.Sleep(1 * time.Second)

    elapsed := time.Since(start)
    fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}
```
