# 使用 defer、panic 和 recover 函式控制

現在讓我們看看一些對 Go 而言獨特的控制流程：defer、panic 和 recover。 這些函式每個都有數個使用案例。 您將在這裡探索最重要的使用案例。

## Defer 函式
在 Go 中，defer 陳述式會延後執行函式 (包括任何參數)，直到包含 defer 陳述式的函式結束為止。 一般而言，當您想要避免忘記關閉檔案或執行清除程序等工作時，就會延遲函式。

您可以延遲任意數量的函式。 defer 陳述式會反向執行，從最後一個到第一個。

執行下列範例程式碼以查看此模式的運作方式：

```Go
package main

import "fmt"

func main() {
    for i := 1; i <= 4; i++ {
        defer fmt.Println("deferred", -i)
        fmt.Println("regular", i)
    }
}
```
以下是程式碼輸出：

輸出
```
regular 1
regular 2
regular 3
regular 4
deferred -4
deferred -3
deferred -2
deferred -1
```

請注意，在此範例中，每次延遲 fmt.Println("deferred", -i) 時，都會儲存 i 的值，並將函式呼叫新增至佇列。 在 main() 函式結束列印 regular 值後，才會執行所有延遲的呼叫。 請注意，延遲呼叫的輸出會以反向順序 (後進先出) 從佇列中取出。

defer 函式的典型使用案例是用完檔案後關閉檔案。 以下是範例：

```Go
package main

import (
    "io"
    "os"
    "fmt"
)

func main() {
    newfile, error := os.Create("learnGo.txt")
    if error != nil {
        fmt.Println("Error: Could not create file.")
        return
    }
    defer newfile.Close()

    if _, error = io.WriteString(newfile, "Learning Go!"); error != nil {
	    fmt.Println("Error: Could not write to file.")
        return
    }

    newfile.Sync()
}
```
在您建立或開啟檔案之後，您要延遲 .Close() 函式以免完成後忘記關閉檔案。

---
## Panic 函式
執行階段錯誤會使 Go 程式異常，例如嘗試使用超出範圍的索引或對 nil 指標的取值來存取陣列。 您也可以強制讓程式異常。

內建 panic() 函式會停止 Go 程式中正常的控制流程。 當您使用 panic 呼叫時，任何延遲的函式呼叫都會正常執行。 這個程序會繼續執行堆疊，直到所有函式都傳回為止。 然後程式當機並顯示一則記錄訊息。 此訊息包含所有錯誤和堆疊追蹤，可協助您診斷問題的根本原因。

當您呼叫 panic() 函式時，您可以新增任何值為引數。 通常，您會傳送說明緊急狀況發生原因的錯誤訊息。

例如，下列程式碼會結合 panic 和 defer 函式。 請嘗試執行此程式碼，以查看控制流程中斷的方式。 請注意，清除流程仍會執行。

```Go

package main

import "fmt"

func highlow(high int, low int) {
    if high < low {
        fmt.Println("Panic!")
        panic("highlow() low greater than high")
    }
    defer fmt.Println("Deferred: highlow(", high, ",", low, ")")
    fmt.Println("Call: highlow(", high, ",", low, ")")

    highlow(high, low + 1)
}

func main() {
    highlow(2, 0)
    fmt.Println("Program finished successfully!")
}
```
輸出如下：

```
Call: highlow( 2 , 0 )
Call: highlow( 2 , 1 )
Call: highlow( 2 , 2 )
Panic!
Deferred: highlow( 2 , 2 )
Deferred: highlow( 2 , 1 )
Deferred: highlow( 2 , 0 )
panic: highlow() low greater than high

goroutine 1 [running]:
main.highlow(0x2, 0x3)
	/tmp/sandbox/prog.go:13 +0x34c
main.highlow(0x2, 0x2)
	/tmp/sandbox/prog.go:18 +0x298
main.highlow(0x2, 0x1)
	/tmp/sandbox/prog.go:18 +0x298
main.highlow(0x2, 0x0)
	/tmp/sandbox/prog.go:18 +0x298
main.main()
	/tmp/sandbox/prog.go:6 +0x37

Program exited: status 2.
```
這是程式碼執行時的狀況：

* 所有作業都正常執行。 程式會列印傳入 highlow() 函式的最高值和最低值。
* 如果 low 的值大於 high 的值，程式會異常。 您會看見 Panic! 訊息。 此時，控制流程會中斷，而所有延遲的函式都會開始列印 Deferred... 訊息。
* 程式當機，而您會看到完整的堆疊追蹤。 您不會看到 Program finished successfully! 訊息。

發生非預期的重大錯誤時，通常會執行呼叫 panic() 函式。 若要避免程式損毀，您可以使用另一個名為 recover() 的函式。

## Recover 函式
有時您可能想要避免程式當機，而改為在內部報告錯誤。 或者，您可能想要先清除這一團亂，再讓程式當機。 例如，您可能想要關閉對某項資源的所有連線，以避免發生更多問題。

Go 提供的內建 recover() 函式，可讓您在發生異常後重新取得控制權。 您只需在呼叫 recover 函式的同時也也呼叫 defer 函式。 如果呼叫 recover() 函式，會傳回 nil，但不影響正常執行。

請嘗試修改之前的程式碼的 main 函式，新增對 recover() 函式的呼叫，如下所示：

```Go
func main() {
    defer func() {
	handler := recover()
        if handler != nil {
            fmt.Println("main(): recover", handler)
        }
    }()

    highlow(2, 0)
    fmt.Println("Program finished successfully!")
}
```
執行程式後的輸出應如下所示：

```
Call: highlow( 2 , 0 )
Call: highlow( 2 , 1 )
Call: highlow( 2 , 2 )
Panic!
Deferred: highlow( 2 , 2 )
Deferred: highlow( 2 , 1 )
Deferred: highlow( 2 , 0 )
main(): recover from panic highlow() low greater than high

Program exited.
```
發現與前一版本的差異了嗎？ 主要的差異是您不會再看到堆疊追蹤錯誤。

在 main() 函式中，您延後了呼叫 recover() 函式的匿名函式。 當程式發生緊急狀況時，呼叫 recover() 不會傳回 nil。 此時，您可以做些事來解決這個問題，但在這個案例中，您只要列印一些東西即可。

panic 和 recover 函式的組合是 Go 處理例外狀況的慣用方式。 其他程式設計語言則使用 try/catch 區塊。 Go 偏好您在這裡探索的方法。

---

[下一個單元](./flow-5.md):練習 - 使用 Go 的控制流程