package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	wantBody := "Hello, World!"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(wantBody))
	}))

	defer server.Close()

	got, err := Fetch(server.URL)
	if err != nil {
		t.Errorf("error getting reference: %v", err)
	}

	if got != wantBody {
		t.Errorf("expected %s, got %s", wantBody, got)
	}
}
