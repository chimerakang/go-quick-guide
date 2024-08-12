# 上下文 context
by [@chimerakang](https://github.com/chimerakang)

---
## 摘要
* 將路由(router)獨立出來，方便之後增強。
* 設計上下文(Context)，封裝Request 和Response ，提供對JSON、HTML 等回傳類型的支援。

---

## 設計Context
### 必要性
* 對Web服務來說，是根據請求`*http.Request`，建構回應`http.ResponseWriter`。但這兩個物件提供的介面太過詳細，例如我們要建構一個完整的回應，需要考慮訊息頭(Header)和訊息體(Body)，而Header 包含了狀態碼(StatusCode)，訊息類型(ContentType)等幾乎每次請求都需要設定的資訊。因此，如果不進行有效的封裝，那麼框架的使用者將需要大量重複，繁雜的程式碼，而且容易出錯。針對常用場景，能夠有效率地建構出HTTP 響應是一個好的框架必須考慮的點。

用返回JSON 資料作比較，感受下封裝前後的差距。

封裝前
```go
obj = map[string]interface{}{
    "name": "gootutu",
    "password": "1234",
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
encoder := json.NewEncoder(w)
if err := encoder.Encode(obj); err != nil {
    http.Error(w, err.Error(), 500)
}
    
```

VS 封裝後：
```go
c.JSON(http.StatusOK, gee.H{
    "username": c.PostForm("username"),
    "password": c.PostForm("password"),
})
```

* 針對使用場景，封裝`*http.Request`和`http.ResponseWriter`的方法，簡化相關介面的調用，只是設計Context 的原因之一。對於框架來說，還需要支撐額外的功能。例如，將來解析動態路由/hello/:name，參數:name的值放在哪呢？再例如，框架需要支援中間件，那中間件產生的資訊放在哪呢？ Context 隨著每個請求的出現而產生，請求的結束而銷毀，和當前請求強相關的資訊都應由Context 承載。因此，設計Context 結構，擴展性和複雜性留在了內部，而對外簡化了介面。路由的處理函數，以及將要實現的中間件，參數都統一使用Context 實例， Context 就像一次會話的百寶箱，可以找到任何東西。

### 具體實現
day2-context/gee/context.go
```go
type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
```
* 程式碼最開頭，給`map[string]interface{}`起了一個別名`goo.H`，建構JSON資料時，顯得更簡潔。
* `Context`目前只包含了`http.ResponseWriter`和`*http.Request`，另外提供了對Method 和Path 這兩個常用屬性的直接存取。
* 提供了存取`Query`和`PostForm`參數的方法。
* 提供了快速建構String/Data/JSON/HTML回應的方法。

---
## 路由(Router)
我們將和路由相關的方法和結構提取了出來，放到了一個新的檔案中`router.go`，方便我們下次對router 的功能進行增強，例如提供動態路由的支援。 router 的handle 方法做了一個細微的調整，即handler 的參數，變成了Context。
```go
package goo

import (
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
	
```

---

## 框架入口
```go
package goo

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
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
	c := newContext(w, req)
	engine.router.handle(c)
}

```
將`router`相關的程式碼獨立後，`goo.go`簡單了不少。最重要的還是透過實作了ServeHTTP 接口，接管了所有的HTTP 請求。相較於之前的程式碼，這個方法也有細微的調整，在呼叫router.handle 之前，建構了一個Context 物件。這個物件目前還非常簡單，只是包裝了原來的兩個參數，之後我們會慢慢地為Context插上翅膀。

## 使用效果
為了展示框架的成果，我們實作另外一個簡短的程式

``` go
package main

import (
	"net/http"
	"goo"
)

func main() {
	r := goo.New()
	r.GET("/", func(c *goo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Goo</h1>")
	})
	r.GET("/hello", func(c *goo.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *goo.Context) {
		c.JSON(http.StatusOK, goo.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}

```
* `Handler`的參數變成成了`goo.Context`，提供了查詢`Query/PostForm`參數的功能。
* `goo.Context`封裝了`HTML/String/JSON`函數，能夠快速建構HTTP響應。

運行`go run web4.go`，借助curl ，一起看成果。
```
$ curl -i http://localhost:9999
HTTP/1.1 200 OK
Content-Type: text/html
Date: Mon, 12 Aug 2024 03:57:49 GMT
Content-Length: 18

<h1>Hello Goo</h1>

$ curl http://localhost:9999/hello?=chimera
hello , you're at /hello

$ curl http://localhost:9999/xxx
404 NOT FOUND: /xxx

```

---
## Next
[網頁-路由](./web-3.md)