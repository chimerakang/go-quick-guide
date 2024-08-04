# 使用 for 運算式對資料進行迴圈

迴圈是另一個常見的控制流程。 Go 只使用一個迴圈建構，即 for 迴圈， 但表示迴圈的方式不只一種。 在這部分的課程中，您將認識 Go 支援的迴圈模式。

## 基本 for 迴圈語法
如同 if 陳述式和 switch 陳述式，for 迴圈運算式也不需要括弧， 但需要有大括弧。

分號 (;) 會分隔 for 迴圈的三個元件：

在第一次反覆運算前執行的初始陳述式 (選擇性)。
在每次反覆運算前評估的條件運算式。 當此條件為 false 時，迴圈就會停止。
每次反覆運算結束時執行的 poststatement (選擇性)。
如您所見，Go 中的 for 迴圈類似 C、Java 和 C# 等程式設計語言中的 for 迴圈，

Go 中最簡單的 for 迴圈看起來像這樣：

```Go
func main() {
    sum := 0
    for i := 1; i <= 100; i++ {
        sum += i
    }
    fmt.Println("sum of 1..100 is", sum)
}
```
讓我們看看在 Go 中撰寫迴圈的其他方式。

## 空白的 prestatement 和 poststatement
在某些程式設計語言中，當只需要條件運算式時，您要使用 while 關鍵字撰寫迴圈模式。 Go 沒有 while 關鍵字， 但您可以改用 for 迴圈，並利用 Go 讓 prestatement 和 poststatement 成為選擇性的事實。

使用下列程式碼片段確認您可在不使用 prestatement 和 poststatement 的情況下使用 for 迴圈。

```Go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    var num int64
    rand.Seed(time.Now().Unix())
    for num != 5 {
        num = rand.Int63n(15)
        fmt.Println(num)
    }
}
```
只要 num 變數不是 5，程式就會列印亂數。

## 無限迴圈和 break 陳述式
無限迴圈是您可以在 Go 中撰寫的另一個迴圈模式， 在這種情況下，您不用撰寫條件運算式或 prestatement 或 poststatement， 可以用自己的方式跳出迴圈； 否則，邏輯將永遠不會結束。 請使用 break 關鍵字，讓邏輯結束迴圈。

若要撰寫適當的無限迴圈，請以大括弧括住 for 關鍵字，如下所示：

```Go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    var num int32
    sec := time.Now().Unix()
    rand.Seed(sec)

    for {
        fmt.Print("Writing inside the loop...")
        if num = rand.Int31n(10); num == 5 {
            fmt.Println("finish!")
            break
        }
        fmt.Println(num)
    }
}
```
每次執行此程式碼時，都會得到不同的輸出。

## Continue 陳述式
在 Go 中，您可以使用 continue 關鍵字略過迴圈目前的反覆項目。 例如，您可以使用這個關鍵字先執行驗證，再繼續執行迴圈。 或者在撰寫無限迴圈，需要等待資源變成可用時，使用此關鍵字。

此範例採用 continue 關鍵字：

```Go
package main

import "fmt"

func main() {
    sum := 0
    for num := 1; num <= 100; num++ {
        if num%5 == 0 {
            continue
        }
        sum += num
    }
    fmt.Println("The sum of 1 to 100, but excluding numbers divisible by 5, is", sum)
}
```

這個範例有一個從 1 到 100 反覆運算的 for 迴圈，會將目前的數字加入每個反覆項目的總和。 在迴圈目前的反覆運算中，略過每個可由 5 整除的數字，不加入至總和。

---
[下一個單元](./flow-4.md): 使用 defer、panic 和 recover 函式控制