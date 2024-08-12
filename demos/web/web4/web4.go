package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]
// curl http://localhost:9999/world
// 404 NOT FOUND: /world

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
		// expect /hello?name=gootutu
		c.String(http.StatusOK, "hello %s, password:%s,  you're at %s\n", c.Query("name"), c.Query("pwd"), c.Path)
	})

	r.POST("/login", func(c *goo.Context) {
		c.JSON(http.StatusOK, goo.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
