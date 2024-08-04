# 流程控制

## 使用 if/else 運算式測試條件

if/else 陳述式是所有程式設計語言中最基本的控制流程， 所以在 Go 中，if/else 陳述式很簡單， 但是，您必須先了解幾項差異，才會開始習慣撰寫 Go 程式。

讓我們看看 if 陳述式的 Go 語法。

if 陳述式的語法
與其他程式設計語言不同，Go 中的「條件不需要括弧」， else 子句是選擇性子句。 但仍然需要大括弧。 而且，Go 不支援使用三元 if 陳述式縮減行數，所以每次都必須撰寫完整的 if 陳述式。

以下是 if 陳述式的基本範例：

```go
package main

import "fmt"

func main() {
    x := 27
    if x%2 == 0 {
        fmt.Println(x, "is even")
    }
}
```

在 VS Code 中，如果您的 Go 語法條件包含括弧，則系統會在儲存程式時自動移除括弧。

## 複合 if 陳述式
Go 支援複合 if 陳述式， 您可以使用 else if 陳述式語句建立巢狀陳述式。 以下是範例：

```Go
package main

import "fmt"

func givemeanumber() int {
    return -1
}

func main() {
    if num := givemeanumber(); num < 0 {
        fmt.Println(num, "is negative")
    } else if num < 10 {
        fmt.Println(num, "has only one digit")
    } else {
        fmt.Println(num, "has multiple digits")
    }
}
```

請注意，在這段程式碼中，num 變數會儲存 givemeanumber() 函式傳回的值，而「所有 if 分支都可以使用」此變數。 但是，如果您嘗試在 if 區塊外列印 num 變數的值，則會收到錯誤:

```Go
package main

import "fmt"

func somenumber() int {
    return -7
}
func main() {
    if num := somenumber(); num < 0 {
        fmt.Println(num, "is negative")
    } else if num < 10 {
        fmt.Println(num, "has 1 digit")
    } else {
        fmt.Println(num, "has multiple digits")
    }

    fmt.Println(num)
}
```
當您執行程式時，錯誤輸出看起來像這樣：

輸出

```
# command-line-arguments
./main.go:17:14: undefined: num
```

在 Go 中，宣告 if 區塊內的變數是慣用的。 使用慣例有效地進行程式設計在 Go 中很常見。

---

## Next
[流程控制2](./flow-2.md)
