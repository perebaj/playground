package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpperCaseHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil) // new incoming request
	w := httptest.NewRecorder()                                        // records it's mutations for later inspection
	upperCaseHandler(w, req)                                           // call the handler
	res := w.Result()                                                  // get the result
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(data))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "ABC" {
		t.Errorf("expected ABC got %v", string(data))
	}
}
