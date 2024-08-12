# 中間件 Middleware
by [@chimerakang](https://github.com/chimerakang)

---
## 摘要
* 設計並實作Web 框架的中間件(Middlewares)機制。
* 實作通用的`Logger`中間件，能夠記錄請求到回應所花費的時間

---
## 中間件是什麼
中間件(middlewares)，簡單說，就是非業務的技術類別元件。 Web 框架本身不可能去理解所有的業務，因而不可能實現所有的功能。因此，框架需要有一個介面，讓使用者自己定義功能，嵌入框架中，彷彿這個功能是框架原生支援的一樣。因此，對中間件而言，需要考慮2個比較關鍵的點：

* 插入點在哪？使用框架的人並不關心底層邏輯的具體實現，如果插入點太底層，中間件邏輯就會非常複雜。如果插入點離使用者太近，那麼和使用者直接定義一組函數，每次在Handler 中手動呼叫沒有多大的優勢了。
* 中間件的輸入是什麼？中間件的輸入，決定了擴充能力。暴露的參數太少，使用者發揮空間有限。

那對於一個Web 框架而言，中間件該設計成什麼樣子呢？

---
## 中間件設計
Goo 的中間件的定義與路由對映的Handler 一致，處理的輸入是`Context`物件。插入點是框架接收到請求初始化`Context`物件後，允許使用者使用自己定義的中間件做一些額外的處理，例如記錄日誌等，以及對進行`Context`二次加工。另外透過呼叫`(*Context).Next()`函數，中間件可等待使用者自己定義的`Handler`處理結束後，再做一些額外的操作，例如計算本次處理所用時間等。即`Goo` 的中間件支援用戶在請求被處理的前後，做一些額外的操作。舉個例子，我們希望最終能夠支援如下定義的中間件，`c.Next()`表示等待執行其他的中間件或使用者的`Handler`：

### logger.go
```go
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
```	 

另外，支援設定多個中間件，依序進行呼叫。

我們上一篇文章分組控制[Group Control](./web-4.md)講到，中間件是應用在`RouterGroup`上的，應用在最頂層的Group，相當於作用於全局，所有的請求都會被中間件處理。那為什麼不作用在每一條路由規則上呢？作用在某路由規則，那不如使用者直接在Handler 中呼叫直覺。只作用在某路由規則的功能通用性太差，不適合定義為中介軟體。

我們之前的框架設計是這樣的，當接收到請求後，匹配路由，該請求的所有資訊都保存在`Context`中。中間件也不例外，接收到請求後，應查找所有應作用於該路由的中間件，保存在`Context`中，依序進行呼叫。為什麼依序呼叫後，還需要在`Context`中儲存呢？因為在設計中，中間件不僅作用在處理流程前，也可以作用在處理流程後，也就是在使用者定義的`Handler` 處理完畢後，還可以執行剩下的操作。

為此，我們為Context增加了2個參數，定義了`Next`方法：

### context.go
```go
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
```

`index`是記錄目前執行到第幾個中間件，當在中間件中呼叫`Next`方法時，控制權交給了下一個中間件，直到呼叫到最後一個中間件，然後再從後往前，呼叫每個中間件在`Next`方法之後定義的部分。如果我們將使用者在映射路由時定義的`Handler`新增到`c.handlers`清單中，結果會怎麼樣呢？想必你已經猜到了。
```go
func A(c *Context) {
    part1
    c.Next()
    part2
}
func B(c *Context) {
    part3
    c.Next()
    part4
}
```
假設我們應用了中間件A 和B，和路由對映的Handler。`c.handlers`是這樣的[A, B, Handler]，`c.index`初始化為-1。調用`c.Next()`，接下來的流程是這樣的：

* c.index++，c.index 變成 0
* 0 < 3，呼叫c.handlers[0]，即 A
* 執行part1，呼叫c.Next()
* c.index++，c.index 變成 1
* 1 < 3，呼叫c.handlers[1]，即 B
* 執行part3，呼叫c.Next()
* c.index++，c.index 變成 2
* 2 < 3，呼叫c.handlers[2]，即Handler
* Handler 調用完畢，回到B 中的part4，執行part4
* part4 執行完畢，返回A 中的part2，執行part2
* part2 執行完畢，結束。

一句話說清楚重點，最終的順序是`part1 -> part3 -> Handler -> part 4 -> part2`。恰恰滿足了我們對中間件的要求，接下來看調用部分的程式碼，就能全部串起來了。

程式碼實現
定義Use函數，將中間件套用到某個Group 。
### 修改goo.go
```go
// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
```	

ServeHTTP 函數也有變化，當我們接收到一個具體請求時，要判斷該請求適用於哪些中間件，在這裡我們簡單透過URL 的前綴來判斷。得到中間件清單後，賦值給`c.handlers`。

* handle 函數中，將從路由匹配得到的Handler 加入`c.handlers`清單中，執行`c.Next()`。

### 修改 router.go
``` go
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
```

### 測試
```
func onlyForV2() goo.HandlerFunc {
	return func(c *goo.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := goo.New()
	r.Use(goo.Logger()) // global midlleware
	r.GET("/", func(c *goo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello goo</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *goo.Context) {
			// expect /hello/gooktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
```

`goo.Logger()`也就是我們一開始就介紹的中間件，我們將這個中間件和框架程式碼放在了一起，作為框架預設提供的中間件。在這個例子中，我們將`goo.Logger()`應用在了全局，所有的路由都會應用該中間件。`onlyForV2()`是用來測試功能的，僅在v2對應的Group 中應用了。

接下來使用curl 測試，可以看到，v2 Group 2個中間件都生效了。
```
$ curl http://localhost:9999/
>>> log
2019/08/01 01:37:38 [200] / in 3.14µs

(2) global + group middleware
$ curl http://localhost:9999/v2/hello/gootutu
>>> log
2024/08/01 01:38:48 [200] /v2/hello/gootutu in 61.467µs for group v2
2024/08/01 01:38:48 [200] /v2/hello/gootutu in 281µs
```

---
## Next
[網頁-模板](./web-6.md)