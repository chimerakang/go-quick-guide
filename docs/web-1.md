# 網頁基礎 web
by [@chimerakang](https://github.com/chimerakang)

---
## 介紹
由於Golang優秀的併發處理，很多公司使用Golang編寫微服務。 對於Golang來說，只需要短短幾行代碼就可以實現一個簡單的Http伺服器。 加上Golang的協程，這個伺服器可以擁有極高的性能。 然而，正是因為代碼過於簡單，我們才應該去研究他的底層實現，也才知道怎麼使用。

[參考:Web工作方式](https://willh.gitbook.io/build-web-application-with-golang-zhtw/03.0/03.1)

在本文中，會以自頂向下的方式，從如何使用，到如何實現，一點點的分析Golang中`net/HTTP`這個包中關於Http伺服器的實現方式。

## 標準函式庫啟動Web服務
Go語言內建了net/http庫，封裝了HTTP網路程式設計的基礎的接口，我們實現的`GooWeb` 框架是基於`net/http`的。我們接下來透過一個例子，簡單介紹下這個函式庫的使用。

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q \n", k, v)
	}
}

```


我們設定了2個路由，`/`和`/hello`，分別綁定`indexHandler`和`helloHandler`， 根據不同的HTTP請求會呼叫不同的處理函數。訪問`/`，回應是`URL.Path = /`，而`/hello`的回應則是請求頭(header)中的鍵值對訊息。

用curl 這個工具測試一下，將會得到如下的結果。

```
$ curl http://localhost:9999/ 
URL.Path = "/"

$ curl http://localhost:9999/hello
Header["User-Agent"] = ["curl/7.68.0"] 
Header["Accept"] = ["*/*"]
```
main函數的最後一行，是用來啟動Web 服務的，第一個參數是位址，`:9999`表示在9999埠監聽。而第二個參數則代表處理所有的HTTP請求的實例，`nil`代表使用標準庫中的實例處理。第二個參數，則是我們基於`net/http`標準函式庫實現Web框架的入口。

## 實作http.Handler介面
``` go
package http

type Handler interface { 
    ServeHTTP(w ResponseWriter, r *Request) 
}

func ListenAndServe(address string , h Handler) error
```

第二個參數的型別是什麼呢？透過查看`net/http`的源碼可以發現，`Handler`是一個接口，需要實現方法`ServeHTTP`，也就是說，只要傳入任何實現了`ServerHTTP`接口的實例，所有的HTTP請求，就都交給了該實例處理了。馬上來試試吧。

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}

```

* 我們定義了一個空的結構體`Engine`，實作了方法`ServeHTTP`。這個方法有2個參數，第二個參數是`Request`，該物件包含了該HTTP請求的所有的信息，例如請求位址、Header和Body等資訊；第一個參數是ResponseWriter，利用ResponseWriter可以建構針對該請求的響應。

* 在main函數中，我們給ListenAndServe方法的第二個參數傳入了剛才建立的engine實例。至此，我們走出了實現Web框架的第一步，即，將所有的HTTP請求轉向了我們自己的處理邏輯。還記得嗎，在實作Engine之前，我們呼叫http.HandleFunc實作了路由和Handler的映射，也就是只能針對特定的路由寫入處理邏輯。比如/hello。但是在實現Engine之後，我們攔截了所有的HTTP請求，擁有了統一的控制入口。在這裡我們可以自由定義路由映射的規則，也可以統一加入一些處理邏輯，例如日誌、例外處理等。

* 程式碼的運行結果與之前的是一致的。

## 框架的雛形
我們接下來重新組織上面的程式碼，搭建出整個框架的雛形。

最終的程式碼目錄結構是這樣的。
```
goo/ 
  |--goo.go 
  |--go.mod 
web3.go 
go.mod
```
go.mod 的內容
```
module example

go 1.21.5

require goo v0.0.0

replace goo => ./goo
```
在`go.mod`中使用replace將goo 指向`./goo`
從`go 1.11` 版本開始，引用相對路徑的package 需要上述使用方式。
```go

package main

import (
	"fmt"
	"net/http"

	"goo"
)

func main() {
	r := goo.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})


	r.Run(":9999")
}

```
看到這裡，使用New()建立`goo` 的實例，使用`GET()`方法新增路由，最後使用`Run()`啟動`Web`服務。這裡的路由，只是靜態路由，不支援`/hello/:name`這樣的動態路由，動態路由我們會在下次實作。

### goo.go

```go
package goo

import (
	"fmt"
	"net/http"
)

// HandlerFunc defines the request handler used by goo
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of goo.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
```


那麼goo.go就是重頭戲了。我們重點介紹一下這部分的實作。

* 首先定義了類型`HandlerFunc`，這是提供給框架使用者的，用來定義路由映射的處理方法。我們在`Engine`中，新增了一張路由映射表`router`，key 由請求方法和靜態路由位址構成，例如`GET-/`、`GET-/hello`、`POST-/hello`，這樣針對相同的路由，如果請求方法不同,可以映射不同的處理方法(Handler)，value 是用戶映射的處理方法。

* 當使用者呼叫`(*Engine).GET()`方法時，會將路由和處理方法註冊到映射表router中，`(*Engine).Run()`方法，是`ListenAndServe`的包裝。

* `Engine`實現的`ServeHTTP`方法的功能是，解析請求的路徑，尋找路由映射表，如果查到，就執行註冊的處理方法。如果檢查不到，就回傳`404 NOT FOUND`。

執行go run main.go，再用curl工具訪問，結果與最開始的一致。
```
$ curl http://localhost:9999/ 
URL.Path = "/"
 $ curl http://localhost:9999/hello 
Header[ "Accept" ] = [ "*/*" ] 
Header[ "User-Agent" ] = [ "curl/7.54.0" ] 
curl http://localhost:9999/world 
404 NOT FOUND: /world
```

至此，整個框架的原型已經出來了。實作了路由映射表，提供了使用者註冊靜態路由的方法，包裝了啟動服務的函數。當然，到目前為止，我們還沒有實現比net/http標準庫更強大的能力，不用擔心，很快就可以將動態路由、中間件等功能添加上去了。

---
## Next
[網頁-上下文](./web-2.md)