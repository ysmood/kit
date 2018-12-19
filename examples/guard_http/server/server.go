package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
	})

	fmt.Println("listen on 8080")

	http.ListenAndServe(":8080", nil)
}
