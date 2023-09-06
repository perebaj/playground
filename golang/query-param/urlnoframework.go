package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

// Curl examples:
// curl -X GET "http://localhost:8080/unescape?name=John%20Doe"
// "http://localhost:8080/unescape?name=jonathan+adads"
// "http://localhost:8080/unescape?name=Robert%3DWilliams"

// curl -X GET "http://localhost:8080/nounescape?name=John%20Doe"
// "http://localhost:8080/nounescape?name=jonathan+adads"
// "http://localhost:8080/nounescape?name=Robert%3DWilliams"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/unescape", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		nameUnescaped, err := url.QueryUnescape(query.Get("name"))
		if err != nil {
			log.Error().Err(err).Msg("unescape")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		nmBytes, err := json.Marshal(nameUnescaped)
		if err != nil {
			log.Error().Err(err).Msg("json.Marshal")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(nmBytes))
	})

	mux.HandleFunc("/nounescape", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		name := query.Get("name")
		nmBytes, err := json.Marshal(name)
		if err != nil {
			log.Error().Err(err).Msg("json.Marshal")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(nmBytes))
	})

	http.ListenAndServe(":8080", mux)
}
