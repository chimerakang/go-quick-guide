# Buffer Channel 有緩衝通道

如您前面所學，通道根據預設是無法緩衝。 這表示僅當有接收作業時，才會接受傳送作業。 否則將會限制程式一直等待。

有時候，您會需要在 Goroutine 之間進行這種類型的同步。 不過，有時候您可能只想實作並行，不需要限制 Goroutine 彼此相互通訊的方式。

因為可緩衝的通道與佇列的運作方式相同，所以不會限制程式傳送及接收資料。 當您建立通道時，可以限制此佇列的大小，如下所示：

```Go

ch := make(chan string, 10)
```
每當有項目傳送至通道時，就會新增至佇列。 然後，接收作業會從佇列中移除該項目。 當通道滿載時，所有傳送作業都會等候，直到有空間可以保存資料為止。 相反地，若通道是空的，而且有讀取作業，將會限制通道，直到讀取開始為止。

以下提供一簡單範例，協助您了解 buffer channel：

```Go
package main

import (
    "fmt"
)

func send(ch chan string, message string) {
    ch <- message
}

func main() {
    size := 4
    ch := make(chan string, size)
    send(ch, "one")
    send(ch, "two")
    send(ch, "three")
    send(ch, "four")
    fmt.Println("All data sent to the channel ...")

    for i := 0; i < size; i++ {
        fmt.Println(<-ch)
    }

    fmt.Println("Done!")
}
```
當您執行程式時，會看到下列輸出：

輸出

```
All data sent to the channel ...
one
two
three
four
Done!
```
您可能會覺得，我們在這裡沒有做什麼不同的事，您想得沒錯。 但是，讓我們看看，當您將 size 變數變更為較小的數字時 (您也可以試試較大的數字) 會如何，如下所示：

```Go
size := 2
```
當您重新執行程式時，會收到下列錯誤：

輸出

```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.send(...)
        /Users/developer/go/src/concurrency/main.go:8
main.main()
        /Users/developer/go/src/concurrency/main.go:16 +0xf3
exit status 2
```
原因是對 send 函式的呼叫是連續的。 您並沒有建立新的 Goroutine。 因此不會有佇列。

通道與 Goroutine 緊密相連。 若無其他 Goroutine 接收來自通道的資料，整個程式可能會永遠處於受限狀態。 如您所見，這確實會發生。

現在讓我們做件有趣的事！ 我們為最後兩個呼叫建立 Goroutine (先前兩個呼叫已正確排入緩衝)，然後執行 for 迴圈四次。 程式碼如下：

```Go
func main() {
    size := 2
    ch := make(chan string, size)
    send(ch, "one")
    send(ch, "two")
    go send(ch, "three")
    go send(ch, "four")
    fmt.Println("All data sent to the channel ...")

    for i := 0; i < 4; i++ {
        fmt.Println(<-ch)
    }

    fmt.Println("Done!")
}
```
當您執行此程式時，程式會如預期般運作。 當您使用通道時，建議您一律使用 Goroutine。

讓我們測試您建立具有超過所需元素的緩衝通道案例。 我們將使用之前使用的範例來檢查 API 並建立一個大小為 10 的緩衝區通道:

```Go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func main() {
    start := time.Now()

    apis := []string{
        "https://management.azure.com",
        "https://dev.azure.com",
        "https://api.github.com",
        "https://outlook.office.com/",
        "https://api.somewhereintheinternet.com/",
        "https://graph.microsoft.com",
    }

    ch := make(chan string, 10)

    for _, api := range apis {
        go checkAPI(api, ch)
    }

    for i := 0; i < len(apis); i++ {
        fmt.Print(<-ch)
    }

    elapsed := time.Since(start)
    fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}

func checkAPI(api string, ch chan string) {
    _, err := http.Get(api)
    if err != nil {
        ch <- fmt.Sprintf("ERROR: %s is down!\n", api)
        return
    }

    ch <- fmt.Sprintf("SUCCESS: %s is up and running!\n", api)
}
```
當您執行程式時，得到的輸出將會與先前相同。 您可以使用較小或較大的數字來變更通道大小，但程式仍可運作。

---
## Unbuffer Channel 與 Buffer Channel
至此，您可能很想知道，兩種通道的使用時機。 這完全取決於您希望 Goroutine 之間通訊的流動方式。 無法緩衝的通道會以同步方式通訊， 以確保每次傳送資料時，程式都會受到限制，直到有人從通道讀取資料為止。

反之，可緩衝的通道會分隔傳送與接收作業， 而不會限制程式，但請務必小心，這可能會造成死結 (如您先前所見)。 當您使用無法緩衝的通道時，可以控制可並行的 Goroutine 數量。 例如，您可能會想要呼叫 API，同時控制每秒執行的呼叫數。 否則，您可能會受到限制。

---
## 通道方向
Go 中的頻道有另一個有趣的功能。 當您在函式的參數中設定通道時，可以指定通道傳送或接收資料。 隨著程式愈來愈大，函式量可能會太多，建議記錄每個通道的意圖，以適當地使用。 若您正在撰寫程式庫，想要公開唯讀的通道，以維持資料一致性。

若要定義通道的方向，可以使用類似於讀取或接收資料的方式。 當您在函式參數中宣告通道時，就可以這麼做。 在函式的參數中定義通道類型的語法：

```Go
chan<- int // it's a channel to only send data
<-chan int // it's a channel to only receive data
```
當您透過通道傳送資料，而該通道只能接收時，您會在編譯程式時收到錯誤。

我們將在下列範例程式中使用兩個函式，各用於讀取資料及傳送資料：

```Go
package main

import "fmt"

func send(ch chan<- string, message string) {
    fmt.Printf("Sending: %#v\n", message)
    ch <- message
}

func read(ch <-chan string) {
    fmt.Printf("Receiving: %#v\n", <-ch)
}

func main() {
    ch := make(chan string, 1)
    send(ch, "Hello World!")
    read(ch)
}
```
當您執行程式時，會看到下列輸出：

輸出
```
Sending: "Hello World!"
Receiving: "Hello World!"
```
此程式會釐清每個函式中各個通道的意圖。 若您嘗試使用「只能接收」的通道來傳送資料，將會收到編譯錯誤。 例如，您可以嘗試執行下列動作：

```Go
func read(ch <-chan string) {
    fmt.Printf("Receiving: %#v\n", <-ch)
    ch <- "Bye!"
}
```
當您執行程式時，會看到下列錯誤：

輸出

```
# command-line-arguments
./main.go:12:5: invalid operation: ch <- "Bye!" (send to receive-only type <-chan string)
```
比起誤用通道，收到編譯錯誤是比較理想的情況。

---
## 多工
最後，我們來看看如何使用 select 關鍵字，同時與多個通道互動。 有時候，當您使用多個通道時，您會想要等候事件發生。 比方說，您可能會加入取消作業的邏輯，在程式處理的資料發生異常狀況時使用。

select 陳述式的運作方式就像是 switch 陳述式，但會用於通道。 該陳述式會限制程式執行，直到收到要處理的事件為止。 若陳述式收到多個事件，將會隨機選擇一個。

select 陳述式的特性是會在處理完事件之後結束執行。 若您想要等候更多事件發生，可能需要使用迴圈。

我們使用下列程式看看 select 的實際運作情況：

```Go

package main

import (
    "fmt"
    "time"
)

func process(ch chan string) {
    time.Sleep(3 * time.Second)
    ch <- "Done processing!"
}

func replicate(ch chan string) {
    time.Sleep(1 * time.Second)
    ch <- "Done replicating!"
}

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    go process(ch1)
    go replicate(ch2)

    for i := 0; i < 2; i++ {
        select {
        case process := <-ch1:
            fmt.Println(process)
        case replicate := <-ch2:
            fmt.Println(replicate)
        }
    }
}
```
當您執行程式時，會看到下列輸出：

輸出
```

Done replicating!
Done processing!
```
請注意，replicate 函式會先完成，這就是為什麼您會先在終端機中看到其輸出。 main 函式具有迴圈，因為 select 陳述式在收到事件時就會結束，但我們仍在等待 process 函式完成。

---
## 超時
有時候會出現 goroutine 阻塞的情況，那麼我們如何避免整個程式進入阻塞的情況呢？我們可以利用 select 來設定超時，透過如下的方式實現：

```go
func main() {
    c := make(chan int)
    o := make(chan bool)
    go func() {
        for {
            select {
                case v := <- c:
                    println(v)
                case <- time.After(5 * time.Second):
                    println("timeout")
                    o <- true
                    break
            }
        }
    }()
    <- o
}
```
---

[挑戰](./channel3.md)