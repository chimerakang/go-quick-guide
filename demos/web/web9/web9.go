package main

import (
	"goo"
	"net/http"
)

type student struct {
	Name string
	Age  int8
}

func main() {
	r := goo.Default()
	r.GET("/", func(c *goo.Context) {
		c.String(http.StatusOK, "Hello gootutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *goo.Context) {
		names := []string{"gootutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
