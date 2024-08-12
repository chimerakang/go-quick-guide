# 分組控制 Gruop
by [@chimerakang](https://github.com/chimerakang)

---
## 摘要
*** 實作路由分組控制(Route Group Control)

---
## 分組的意義
分組控制(Group Control)是Web 框架應提供的基礎功能之一。所謂分組，是指路由的分組。如果沒有路由分組，我們需要針對每一個路由進行控制。但是在真實的業務場景中，往往某一組路由需要相似的處理。例如：

* 以`/post`開頭的路由匿名可存取。
* 以`/admin`開頭的路由需要認證。
* 以`/api`開頭的路由是RESTful 接口，可以對接第三方平台，需要三方平台認證。

大部分情況下的路由分組，是以相同的前綴來區分的。因此，我們今天實現的分組控制也是以前綴來區分，並且支援分組的嵌套。例如`/post`是一個分組，`/post/a`和`/post/b`可以是該分組下的子分組。作用在`/post`分組上的中間件(middleware)，也會作用在子分組，子分組還可以應用自己特有的中間件。

中介軟體可以提供框架無限的擴展能力，應用在分組上，可以使得分組控制的收益更為明顯，而不是共享相同的路由前綴這麼簡單。例如`/admin`的分組，可以應用認證middleware；`/`分組應用日誌中間件，`/`是預設的最頂層的分組，也意味著給所有的路由，即整個框架增加了記錄日誌的能力。

提供擴展能力支援中間件的內容，我們將在下一節中介紹。

---
## 分組嵌套

一個Group 物件需要具備哪些屬性呢？首先是前綴(prefix)，例如`/`，或者`/api`；要支持分組嵌套，那麼需要知道當前分組的父親(parent)是誰；當然了，按照我們一開始的分析，中間件是應用在分組上的，那也需要儲存應用在該分組上的中間件(middlewares)。還記得，我們​​之前調用函數`(*Engine).addRoute()`來映射所有的路由規則和Handler 。如果Group物件需要直接映射路由規則的話，例如我們想在使用框架時，這樣呼叫：

```go
r := goo.New() 
v1 := r.Group( "/v1" ) 
v1.GET( "/" , func (c *goo.Context) { 
	c.HTML(http.StatusOK, "<h1>Hello Goo</h1>" ) 
})
```
那麼Group對象，還需要有存取`Router`的能力，為了方便，我們可以在Group中，保存一個指針，指向`Engine`，整個框架的所有資源都是由`Engine`統一協調的，那麼就可以透過`Engine`間接地存取各種介面了。

所以，最後的Group 的定義是這樣的：

### goo.go
```go
RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}
```
我們也可以進一步地抽象，將`Engine`作為最頂層的分組，也就是說`Engine`擁有`RouterGroup`所有的能力。
```go
Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}
```
那我們就可以將和路由有關的函數，都交給`RouterGroup`實作了。
```go
// New is the constructor of goo.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
```
可以仔細觀察下`addRoute`函數，呼叫了`group.engine.router.addRoute`來實現了路由的映射。由於`Engine`從某種意義上繼承了`RouterGroup`的所有屬性和方法，因為`(*Engine).engine` 是指向自己的。這樣實現，我們既可以像原來一樣添加路由，也可以透過分組添加路由。

---
## 使用Demo
測試框架的Demo就可以這樣寫了：
```go
func main() {
	r := goo.New()
	r.GET("/index", func(c *goo.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *goo.Context) {
			c.HTML(http.StatusOK, "<h1>Hello goo</h1>")
		})

		v1.GET("/hello", func(c *goo.Context) {
			// expect /hello?name=gootutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *goo.Context) {
			// expect /hello/gootutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *goo.Context) {
			c.JSON(http.StatusOK, goo.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
```
通過curl 簡單測試：
```
$ curl http://localhost:9999/v1/hello?name=chimera
hello chimera, you're at /v1/hello

$ curl "http://localhost:9999/v2/hello/gooutu"
hello gootutu, you're at /hello/gootutu
```

---
## Next
[網頁-中間件](./web-5.md)