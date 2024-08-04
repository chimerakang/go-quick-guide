# 練習 - 使用 Go 的控制流程

---
請先嘗試自行解決這些練習。 再比較解決方法和您的結果。 如果不記得重要的細節，您隨時可以複習此課程模組

## 撰寫 FizzBuzz 程式
首先，撰寫可列印數字 1 到 100 的程式，並進行下列變更：

如果數字可由 3 整除，則列印 Fizz。
如果數字可由 5 整除，則列印 Buzz。
如果數字可由 3 和 5 整除，則列印 FizzBuzz。
如不符合以上結果，則列印數字。
嘗試使用 switch 陳述式。

## 尋找質數
撰寫程式以尋找小於 20 的所有質數。 質數是大於1的任何數字，可用自身和 1 整除。 「整除」的意思是，除法運算之後沒有餘數。 就像大部分的程式設計語言一樣，Go 會提供一種方法來檢查除法運算是否產生餘數。 我們可以使用模數 % (百分比符號) 運算子。

在此練習中，您將更新名為 findprimes 的函式，以檢查數字是否為質數。 函式有一個整數引數，並傳回布林值。 函式會透過檢查是否有餘數來測試輸入編號是否為質數。 若數字為質數，則函式回傳 true。

使用下列程式碼片段作為起點，並以正確的語法取代 ?? 的所有執行個體：

```Go
package main

import "fmt"

func findprimes(number int) bool {
	for i := 2; i ?? number; i ?? {
        if number ?? i == ?? {
			return false
        }
    }

	if number ?? {
		return true
	} else {
	    return false
	}
}

func main() {
    fmt.Println("Prime numbers less than 20:")

    for number := ??; number ?? 20; number++ {
        if ?? {
            fmt.Printf("%v ", number)
        }
    }
}
```
此程式會檢查數字 1 到 20，並列印數字 (如果是質數)。 如所述修改範例。

在 main 函式中，對要檢查的所有數字執行迴圈。 檢查最後一個數字之後，請結束迴圈。
呼叫 findprimes 函式以檢查數字。 若函式回傳 true，則列印質數。
在迴圈中 findprimes ，從 2 開始重複，直到計數器大於或等於 number 值為止。
檢查目前計數器的數值是否能將 number 整除。 如果是，請結束迴圈。
當 number 是質數時，會傳回 true，否則傳回 false。
提示：請務必正確處理輸入數字為 1 的案例。

---
## 要求一個數字，如為負數，則出現緊急狀況
撰寫要求使用者輸入一個數字的程式。 以下列程式碼片段為起點：

```Go

package main

import "fmt"

func main() {
    val := 0
    fmt.Print("Enter number: ")
    fmt.Scanf("%d", &val)
    fmt.Println("You entered:", val)
}
```

此程式會要求並列印一個數字。 將範例程式碼修改成：

繼續要求整數數字。 結束迴圈的條件應該是使用者輸入負數時。
當使用者輸入負數時，使程式當機， 然後列印堆疊追蹤錯誤。
當數字為 0 時，列印 0 is neither negative nor positive。 繼續要求數字。
當數字為正數時，列印 You entered: X (X 是輸入的數字)。 繼續要求數字。
暫時忽略使用者可能輸入非整數的可能性。

---
[下一個單元](./flow-6.md):解決方案 - 控制流程練習