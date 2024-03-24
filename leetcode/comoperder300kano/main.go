package main

import (
	"fmt"
	"sort"
)

type DB struct {
	InMemory map[string]interface{}
}

func NewDB() DB {
	return DB{
		InMemory: make(map[string]interface{}),
	}
}

type A struct {
	Value string
}

type Row struct {
	Field string
	Value string
}

type Rows []Row

func main() {
	m := make(map[string]map[string]string)

	m["jonathan"] = make(map[string]string)

	m["jonathan"]["BD"] = "jonathan"
	m["jonathan"]["BE"] = "silva"

	m2 := m["jonathan"]

	var respArray Rows
	for k, v := range m2 {
		r := Row{
			Field: k,
			Value: v,
		}
		respArray = append(respArray, r)
	}

	sort.Slice(respArray, func(i, j int) bool {
		return respArray[i].Field < respArray[j].Field
	})

	var resp string
	for i := 0; i < len(respArray); i++ {
		r := respArray[i]
		if i == len(respArray)-1 {
			resp += r.Field + "(" + r.Value + ")"
		} else {
			resp += r.Field + "(" + r.Value + "),"
		}
	}
	fmt.Println(resp)
	fmt.Println(hasPrefix("jonathan", "jon"))
}

func hasPrefix(s string, prefix string) bool {
	return s[0:len(prefix)] == prefix
}
