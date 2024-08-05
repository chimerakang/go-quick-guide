# Channel 通道
by [@chimerakang](https://github.com/chimerakang)

---
在 Go 中，Channel是 Goroutine 之間的通訊機制。 請記住，Go 的並行處理方法是：「請不要藉由共用記憶體進行通訊；相反地，藉由通訊來共用記憶體」。需要將值從一個 Goroutine 傳送到另一個 Goroutine 時，您可以使用Channel。 現在讓我們一起了解Channel如何運作，以及如何使用Channel來撰寫並行的 Go 程式。

---

## Channel的語法
因為Channel是收送資料的通訊機制，所以也有類型的分別。 這表示您只能傳送資料給通道支援的類型。 您可以使用關鍵字 chan 做為通道的資料類型，但您也必須指定要通過Channel的資料類型，例如 int 類型。

每當您在函式中宣告通道，或在參數中指定通道時，都必須使用 chan <type>，例如 chan int。 若要建立通道，您需要使用內建的 make() 函式：

```go
ch := make(chan int)
```
Channel可以執行傳送資料及接收資料兩項作業。 若要指定Channel支援的作業類型，必須使用Channel運算子 <-。 此外，在Channel中執行傳送資料及接收資料兩項作業會受到限制。 後文中將會說明其原因。

當您想要表示通道只能傳送資料時，請在通道之後使用 <- 運算子。 當您想要表示通道接收資料時，請在通道之前使用 <- 運算子，如這些範例所示：

```Go
ch <- x // sends (or writes ) x through channel ch
x = <-ch // x receives (or reads) data sent to the channel ch
<-ch // receives data, but the result is discarded
```
另一項可以在Channel中使用的作業，就是關閉通道。 若要關閉通道，請使用內建的 close() 函式：

```Go
close(ch)
```
當您關閉通道時，表示資料不再透過通道傳送。 若您嘗試將資料傳送到已經關閉的通道，將會導致程式發生問題； 若您嘗試從已經關閉的通道接收資料，還是能讀取傳送的所有資料。 之後每次「讀取」都會傳回一個零值。

現在讓我們回到先前建立的程式，使用通道移除睡眠功能。 首先，我們要在 main 函式中建立字串通道，如下所示：

```Go
ch := make(chan string)
```
接著移除睡眠行 time.Sleep(3 * time.Second)。

如此就能使用通道在 Goroutine 之間進行通訊。 我們不打算列印 checkAPI 函式的結果，而要重構我們的程式碼，透過通道傳送該訊息。 若要使用該函式的通道，必須新增該通道作為參數。 checkAPI 函式看起來應像這樣：

```Go
func checkAPI(api string, ch chan string) {
    _, err := http.Get(api)
    if err != nil {
        ch <- fmt.Sprintf("ERROR: %s is down!\n", api)
        return
    }

    ch <- fmt.Sprintf("SUCCESS: %s is up and running!\n", api)
}
```
請注意，我們必須使用 fmt.Sprintf 函式，因為我們不想列印任何文字，只是要跨通道傳送格式化文字。 此外也請注意，我們會在通道變數之後使用 <- 運算子來傳送資料。

現在您必須變更 main 函式來傳送通道變數，並接收資料來列印輸出，如下所示：

```Go
ch := make(chan string)

for _, api := range apis {
    go checkAPI(api, ch)
}

fmt.Print(<-ch)
```
請注意，在通道之前使用 <- 運算子，表示我們想要從通道讀取資料。

當您重新執行程式時，會出現如下的輸出：

輸出

```
ERROR: https://api.somewhereintheinternet.com/ is down!

Done! It took 0.007401217 seconds!
```
即便我們沒有呼叫睡眠函式，程式也照常運作了對吧？ 但我們並未達成我們的目的。 我們建立了 5 個 Goroutine，卻只有一個產生輸出。 在下節中，我們將會探討程式為何如此運作。

---
## Unbuffer Channel 無緩衝通道
根據預設，當您使用 make() 函式建立Channel時，也會建立一條無緩衝的通道。 無法緩衝的通道會等待有人可以接收資料時，才會執行傳送作業。 如先前所述，在通道中執行傳送資料及接收資料兩項作業會受到限制的原因。 此封鎖作業也是上節中，程式在收到第一則訊息後隨即停止的原因。

我們可以從 fmt.Print(<-ch) 限制程式執行開始說起，因為程式會從通道讀取資料，並等候資料到達。 當程式取得資料之後，就會繼續執行下一行，直到程式完成。

那麼剩下的 Goroutine 又會如何？ 這些 Goroutine 仍會繼續執行，但其中無一執行接聽的工作。 加上程式提早完成，導致有一些 Goroutine 無法傳送資料。 我們可以新增另一個 fmt.Print(<-ch) 來證明這點，如下所示：

```Go
ch := make(chan string)

for _, api := range apis {
    go checkAPI(api, ch)
}

fmt.Print(<-ch)
fmt.Print(<-ch)
```
當您重新執行程式時，會出現如下的輸出：

輸出
```
ERROR: https://api.somewhereintheinternet.com/ is down!
SUCCESS: https://api.github.com is up and running!
Done! It took 0.263611711 seconds!
```
請注意，您現在會看到兩個 API 的輸出。 當您繼續新增更多的 fmt.Print(<-ch) 行時，最終將會讀取所有要傳送給通道的資料。 但是，若您嘗試讀取更多資料，但已無資料傳送過來時會如何？ 結果類似下列範例所示：

```Go
ch := make(chan string)

for _, api := range apis {
    go checkAPI(api, ch)
}

fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)
fmt.Print(<-ch)

fmt.Print(<-ch)
```
當您重新執行程式時，會出現如下的輸出：

輸出
```
ERROR: https://api.somewhereintheinternet.com/ is down!
SUCCESS: https://api.github.com is up and running!
SUCCESS: https://management.azure.com is up and running!
SUCCESS: https://graph.microsoft.com is up and running!
SUCCESS: https://outlook.office.com/ is up and running!
SUCCESS: https://dev.azure.com is up and running!
```
程式能夠正常運作，但還沒有完成。 因為最後一行列印行會繼續等待接收資料，所以不會印出。 您必須使用類似 Ctrl+C 的命令來關閉程式。

上一個範例只在證明讀取資料與接收資料都會受到限制。 若要修正此問題，您可以將程式碼變更為 for 迴圈，並只接收您確定要傳送的資料，如此範例所示：

```Go
for i := 0; i < len(apis); i++ {
    fmt.Print(<-ch)
}
```
這是程式的最終版本，以防您的版本出現問題：

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

    ch := make(chan string)

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
當您重新執行程式時，會出現如下的輸出：

輸出
```
ERROR: https://api.somewhereintheinternet.com/ is down!
SUCCESS: https://api.github.com is up and running!
SUCCESS: https://management.azure.com is up and running!
SUCCESS: https://dev.azure.com is up and running!
SUCCESS: https://graph.microsoft.com is up and running!
SUCCESS: https://outlook.office.com/ is up and running!
Done! It took 0.602099714 seconds!
```
程式會執行應執行的動作。 您不再使用睡眠函式，而會使用通道。 另外也請注意，若未使用並行，現在大約 600 毫秒就能完成，而不需要將近 2 秒。

最後，我們可以說無法緩衝的通道會同步傳送及接收作業。 即使使用並行，通訊仍會同步。

下一個單元:[Buffer Channel](./channel2.md) 