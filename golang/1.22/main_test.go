package main

import (
	"net/http/httptest"
	"testing"
)

func TestRouterNotAllowed(t *testing.T) {
	mux := Router()

	if mux == nil {
		t.Error("mux is nil")
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", nil)

	mux.ServeHTTP(w, req)

	// 405; Method Not Allowed
	if w.Code != 405 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	allow := w.Header().Get("Allow")
	if allow != "GET, HEAD" {
		t.Errorf("expected GET, HEAD, got %s", allow)
	}
}

func TestRouter(t *testing.T) {
	mux := Router()

	if mux == nil {
		t.Error("mux is nil")
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if w.Body.String() != "Hello, World!" {
		t.Errorf("expected Hello, World!, got %s", w.Body.String())
	}
}
