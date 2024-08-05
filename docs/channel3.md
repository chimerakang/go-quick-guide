# 挑戰

面對這項挑戰，您必須改善現有的程式，加快其執行速度。 試試看自行撰寫程式，期間可以回顧前文中您練習的範例。 然後，比較您的解答與下個單元中的解答。

Go 中的並行是很複雜的問題，必須勤於練習才能駕輕就熟。 建議您可以自我挑戰當作練習。

祝您好運！

使用並行快速地計算費氏數值
您可以使用下列程式連續計算費氏數值：

```Go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func fib(number float64) float64 {
    x, y := 1.0, 1.0
    for i := 0; i < int(number); i++ {
        x, y = y, x+y
    }

    r := rand.Intn(3)
    time.Sleep(time.Duration(r) * time.Second)

    return x
}

func main() {
    start := time.Now()

    for i := 1; i < 15; i++ {
        n := fib(float64(i))
    fmt.Printf("Fib(%v): %v\n", i, n)
    }

    elapsed := time.Since(start)
    fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}
```
您必須利用這個現有的程式碼建置兩個程式：

可在其中實作並行的改進版本。 就像現在一樣，新版本應在幾秒鐘內就能完成 (不超過 15 秒)。 您應使用可緩衝的通道。

撰寫新版本計算費氏數值，直到使用者使用 fmt.Scanf() 函式在終端機輸入 quit 為止。 若使用者按 Enter 鍵，就應計算新的費氏數值。 換句話說，您將不再會有介於 1 到 10 之間的迴圈。

使用兩個無法緩衝的通道：一個用於計算費氏數值，一個用於等候使用者的 "quit" 訊息。 您必須使用 select 陳述式。

以下是與程式互動的範例：

輸出
```
1

1

2

3

5

8

13
quit
Done calculating Fibonacci!
Done! It took 12.043196415 seconds!
```

下一個單元: [解決方案](./channel4.md)