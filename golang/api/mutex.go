package main

import (
	"log"
	"net/http"
	"time"
)

type timeHandler struct {
	format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

func main() {
	mux := http.NewServeMux()

	// Initialise the timeHandler in exactly the same way we would any normal
	// struct.
	th := timeHandler{format: time.RFC1123}

	// Like the previous example, we use the mux.Handle() fnction to register
	// this with our ServeMux.
	mux.Handle("/time", th)

	log.Print("Listening...")
	http.ListenAndServe(":3000", mux)
}
