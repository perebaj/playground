package main

import (
	"encoding/json"
	"testing"
)

type User struct {
	Name *string `json:"name"`
	Age  int     `json:"age,omitempty"`
}

func TestUser(t *testing.T) {
	user := User{
		Name: nil,
		Age:  0,
	}

	b, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `{"name":null}` {
		t.Errorf("expected %s, got %s", `{"name":null}`, string(b))
	}
}

func TestUserWithOmitEmpty(t *testing.T) {
	type UserWithOmitEmpty struct {
		Name *string `json:"name,omitempty"`
		Age  int     `json:"age,omitempty"`
	}

	user := UserWithOmitEmpty{
		Name: nil,
		Age:  0,
	}

	b, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `{}` {
		t.Errorf("expected %s, got %s", `{}`, string(b))
	}
}
