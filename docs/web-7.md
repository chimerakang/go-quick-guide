# 錯誤處理 Panic Recover
by [@chimerakang](https://github.com/chimerakang)

---
## 摘要
* 實作錯誤處理機制。

## Panic
Go 語言中，比較常見的錯誤處理方法是傳回error，由呼叫者決定後續如何處理。但如果是無法復原的錯誤，可以手動觸發panic，當然如果在程式運作過程中出現了類似陣列越界的錯誤，panic 也會被觸發。 panic 會中止目前執行的程序，退出。

以下是主動觸發的例子：
```go
// hello.go
func main() {
	fmt.Println("before panic")
	panic("crash")
	fmt.Println("after panic")
}
```
```
$ go run hello.go

before panic
panic: crash

goroutine 1 [running]:
main.main()
        ~/go_demo/hello/hello.go:7 +0x95
exit status 2
```
下面是數組越界觸發的panic
```go
// hello.go
func main() {
	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
}
```

```
$ go run hello.go
panic: runtime error: index out of range [4] with length 3
```

---
## defer
panic 會導致程式被中止，但在退出前，會先處理完目前協程上已經defer 的任務，執行完成後再退出。效果類似java 語言的`try...catch`。
```go
// hello.go
func main() {
	defer func() {
		fmt.Println("defer func")
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
}
```
```
$ go run hello.go 
defer func
panic: runtime error: index out of range [4] with length 3
```

可以defer 多個任務，在同一個函數中defer 多個任務，會逆序執行。即先執行最後defer 的任務。

在這裡，defer 的任務執行完成之後，panic 也會繼續被拋出，導致程式非正常結束。


## recover
Go 語言也提供了recover 函數，可以避免因為panic 發生而導致整個程式終止，recover 函數只在defer 中生效。

```
// hello.go
func test_recover() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after panic")
}

func main() {
	test_recover()
	fmt.Println("after recover")
}
```
```
$ go run hello.go 
defer func
recover success
after recover
```
我們可以看到，recover 捕獲了panic，程式正常結束。test_recover()中的after panic沒有列印，這是正確的，當panic 被觸發時，控制權就交給了defer 。就像在java 中，`try`程式碼區塊中發生了異常，控制權交給了`catch`，接下來執行`catch `程式碼區塊中的程式碼。而在`main()`中列印了`after recover`，說明程式已經恢復正常，繼續往下執行直到結束。

---
## Goo 的錯誤處理機制
對一個Web 框架而言，錯誤處理機制是非常必要的。可能是框架本身沒有完備的測試，導致在某些情況下出現空指標異常等情況。也有可能使用者不正確的參數，觸發了某些異常，例如數組越界，空指標等。如果因為這些原因導致系統宕機，必然是不可接受的。

我們在[之前](./web-6)實現的框架並沒有加入異常處理機制，如果程式碼中存在會觸發panic 的BUG，很容易宕掉。

例如下面的程式碼：
```go
func main() {
	r := goo.New()
	r.GET("/panic", func(c *goo.Context) {
		names := []string{"gooktutu"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}
```

在上面的程式碼中，我們為goo 註冊了路由`/panic`，而這個路由的處理函數內部存在數組越界`names[100]`，如果訪問`localhost:9999/panic`，Web 服務就會當掉。

今天，我們將在goo 中新增一個非常簡單的錯誤處理機制，即在此類錯誤發生時，向使用者返回Internal Server Error，並且在日誌中列印必要的錯誤訊息，方便進行錯誤定位。

我們先前實作了中介軟體機制，錯誤處理也可以作為一個中間件，增強`goo` 框架的能力。

新增檔案`goo/recovery.go`，在這個檔案中實作中間件`Recovery`。
```go
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
```

`Recovery`的實作非常簡單，使用`defer` 掛載上錯誤恢復的函數，在這個函數中呼叫`*recover()*`，捕獲panic，並且將堆疊資訊列印在日誌中，向使用者傳回*Internal Server Error*。

你可能注意到，這裡有一個`trace()`函數，這個函數是用來獲取觸發panic 的堆疊信息，完整程式碼如下：

### recovery.go 
```go
package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
```
在`trace()`中，呼叫了`runtime.Callers(3, pcs[:])`，Callers 用來返回呼叫棧的程式計數器, 第0 個Caller 是Callers 本身，第1 個是上一層trace，第2 個是再上一層的`defer func`。因此，為了日誌簡潔一點，我們跳過了前3 個Caller。

接下來，透過`runtime.FuncForPC(pc)`取得對應的函數，在透過`fn.FileLine(pc)`取得到呼叫函數的檔案名稱和行號，列印在日誌中。

至此，goo 框架的錯誤處理機制就完成了。

### 測試
```go
package main

import (
	"net/http"

	"gee"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello gootutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"gootutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
```
接下來進行測試，先造訪首頁，造訪一個有BUG的`/panic`，服務正常返回。接下來我們再一次成功造訪了主頁，說明服務完全運作正常。

```
1 
2 
3 
4 
5 
6
$curl “http://localhost:9999”
 Hello gooktutu 
$curl “http://localhost:9999/panic”
 { “message”：“內部伺服器錯誤” } 
$curl “http://localhost:9999”
 Hello gooktutu
 ```
我們可以在後台日誌中看到如下內容，引發錯誤的原因和堆疊資訊都被印了出來，透過日誌，我們可以很容易地知道，在`web9.go`的21行的地方出現了`index out of range`錯誤。
```
$ go run web9.go 
2024/08/10 14:07:40 Route  GET - /
2024/08/10 14:07:40 Route  GET - /panic
2024/08/10 14:08:04 [200] / in 0s
2024/08/10 14:08:14 runtime error: index out of range [100] with length 1
Traceback:
        C:/Program Files/Go/src/runtime/panic.go:914
        C:/Program Files/Go/src/runtime/panic.go:114
        G:/GoProjects/src/go-quick-guide/demos/web/web9/web9.go:21
        G:/GoProjects/src/go-quick-guide/demos/web/web9/goo/context.go:41
        G:/GoProjects/src/go-quick-guide/demos/web/web9/goo/recover.go:37
        G:/GoProjects/src/go-quick-guide/demos/web/web9/goo/context.go:41
        G:/GoProjects/src/go-quick-guide/demos/web/web9/goo/logger.go:15
        G:/GoProjects/src/go-quick-guide/demos/web/web9/goo/context.go:41
        G:/GoProjects/src/go-quick-guide/demos/web/web9/goo/router.go:99
        G:/GoProjects/src/go-quick-guide/demos/web/web9/goo/goo.go:130
        C:/Program Files/Go/src/net/http/server.go:2939
        C:/Program Files/Go/src/net/http/server.go:2010
        C:/Program Files/Go/src/runtime/asm_amd64.s:1651

2024/08/10 14:08:14 [500] /panic in 2.9078ms
```

---

參考
[Package runtime](https://pkg.go.dev/runtime)
[Is it possible get information about caller function in Golang?](https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang)