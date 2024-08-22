package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		content := fmt.Sprintf("hello, %s", name)
		fmt.Fprint(w, content)
	})

	http.ListenAndServe(":9999", nil)
}
