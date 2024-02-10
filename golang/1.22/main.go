package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := Router()

	fmt.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func IdHandler(w http.ResponseWriter, r *http.Request) {
	v := r.PathValue("id")
	fmt.Fprintf(w, "id: %s", v)
}

// http patterns: https://pkg.go.dev/net/http#hdr-Patterns
func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", HelloHandler)
	mux.HandleFunc("GET /id/{id}", IdHandler)
	return mux
}
