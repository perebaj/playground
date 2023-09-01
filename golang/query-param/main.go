package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		ids := query["id"]

		idsBytes, err := json.Marshal(ids)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(idsBytes))
	})

	http.ListenAndServe(":8080", r)
}
