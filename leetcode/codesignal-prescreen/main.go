package main

import "fmt"

type DB struct {
	InMemoryDB map[string]map[string]string
}

func NewDb() DB {
	return DB{
		InMemoryDB: make(map[string]map[string]string),
	}
}

type Field struct {
	Field string
	Value string
}

func (db *DB) Set(key string, field string, value string) string {
	_, ok := db.InMemoryDB[key][field]
	if !ok {
		db.InMemoryDB[key] = make(map[string]string)
		db.InMemoryDB[key][field] = value
	} else {
		db.InMemoryDB[key][field] = value
	}
	fmt.Println(db.InMemoryDB)
	return ""
}

func (db *DB) Get(key string, field string) string {
	value, ok := db.InMemoryDB[key][field]
	if !ok {
		return ""
	}
	return value
}

func (db *DB) Delete(key string, field string) string {
	_, ok := db.InMemoryDB[key][field]
	if !ok {
		return "false"
	}

	delete(db.InMemoryDB[key], field)
	return "true"
}

func main() {
	// db := NewDb()
	// db.Set("key", "field", "value")
	// db.Get("key", "field")
	// db.Set("key", "field", "value2")
	// db.Set("key", "field2", "value4")
	// fmt.Println(db.InMemoryDB)
	m := make(map[string]I)

	for i := 0; i < count; i++ {

	}

	i1 := I{
		Key:   "key",
		Field: "field",
	}

	i2 := I{
		Key:   "key2",
		Field: "field2",
	}

	m["key"] = i1
	m["key2"] = i2
	fmt.Println(m)

	m["key"] = i2
}

/*
SCAN <key> — should return a string representing the fields of a record associated with key. The returned string should be in the following format "<field1>(<value1>), <field2>(<value2>), ...", where fields are sorted lexicographically. If the specified record does not exist, returns an empty string.

SCAN_BY_PREFIX <key> <prefix> — should return a string representing some fields of a record associated with key. Specifically, only fields that start with prefix should be included. The returned string should be in the same format as in the SCAN operation with fields sorted in lexicographical order.
asdadsa


Level 2
The database should support displaying data based on filters. Introduce an operation to support printing some fields of a record.

SCAN <key> — should return a string representing the fields of a record associated with key. The returned string should be in the following format "<field1>(<value1>), <field2>(<value2>), ...", where fields are sorted lexicographically. If the specified record does not exist, returns an empty string.

SCAN_BY_PREFIX <key> <prefix> — should return a string representing some fields of a record associated with key. Specifically, only fields that start with prefix should be included. The returned string should be in the same format as in the SCAN operation with fields sorted in lexicographical order.

Examples
The example below shows how these operations should work (the section is scrollable to the right):

Queries	Explanations
queries = [
  ["SET", "A", "BC", "E"],
  ["SET", "A", "BD", "F"],
  ["SET", "A", "C", "G"],
  ["SCAN_BY_PREFIX", "A", "B"],
  ["SCAN", "A"],
  ["SCAN_BY_PREFIX", "B", "B"]
]

returns ""; database state: {"A": {"BC": "E"}}
returns ""; database state: {"A": {"BC": "E", "BD": "F"}}
returns ""; database state: {"A": {"BC": "E", "BD": "F", "C": "G"}}
returns "BC(E), BD(F)"
returns "BC(E), BD(F), C(G)"
returns ""

the output should be ["", "", "", "BC(E), BD(F)", "BC(E), BD(F), C(G)", ""].

Input/Output
*/


/*
func solution(queries [][]string) []string {
    nDB := NewDb()
    var res []string
    for _, v := range queries {

        switch v[0] {
        case "SET":
            resSet := nDB.Set(v[1], v[2], v[3])
            res = append(res, resSet)
        case "GET":
            resGet := nDB.Get(v[1], v[2])
            res = append(res, resGet)
        case "DELETE":
            resDelete := nDB.Delete(v[1], v[2])
            res = append(res, resDelete)
        case "SCAN":
            scan := nDB.Scan(v[1])
            res = append(res, scan)
        case "SCAN_BY_PREFIX":
            respScan := nDB.ScanByPrefix(v[1], v[2])
            res = append(res, respScan)
        }
	}
    return res
}


type DB struct {
    InMemoryDB map[string]map[string]string
}

func NewDb () DB {
    return DB{
        InMemoryDB: make(map[string]map[string]string),
    }
}

type Field struct {
    Field string
    Value string
}

func (db *DB) Set(key string, field string, value string) string{
    _, ok := db.InMemoryDB[key]
    if !ok {
        db.InMemoryDB[key] = make(map[string]string)
        db.InMemoryDB[key][field] = value
    } else {
        db.InMemoryDB[key][field] = value
    }
    return ""
}

func (db *DB) Get(key string, field string) string {
    value, ok := db.InMemoryDB[key][field]
    if !ok {
        return ""
    }
    return value
}

func (db *DB) Delete(key string, field string) string {
    _, ok := db.InMemoryDB[key][field]
    if !ok {
        return "false"
    }

    delete(db.InMemoryDB[key], field)
    return "true"
}

func (db *DB) Scan(key string) string {
    mResult := db.InMemoryDB[key]
    fmt.Println(mResult)
    var results []string
    for k, v := range mResult {
        results = append(results, k + "(" + v + ")")
    }

    fmt.Println(results)
    var res string
    for i := 0; i < len(results); i++ {
        fmt.Println(results[i])
		if i == len(results) - 1{
            res += results[i]
        } else {
            res += results[i] + ","
        }
	}
    return res
}

func (db *DB) ScanByPrefix(key string, prefix string) string {
    mResult := db.InMemoryDB[key]
    var res []string
    for k, _ := range mResult {
        if HasPrefix(k, prefix) {
            res = append(res, k)
        }
    }
    return ""
}

func HasPrefix(s, prefix string) bool {
    return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

/*
[["SET","user1","firstName","Greg"],
 ["SET","user1","lastName","Wright"],
 ["SET","lastName","user1","error"],
 ["GET","user1","lastName"],
 ["GET","user1","firstName"],  erro aqui Greg ""
 ["GET","user1","Greg"],
 ["GET","user1","Wright"],
 ["GET","user1","city"],
 ["DELETE","user1","city"],
 ["SET","user1","city","London"],
 ["GET","user1","firstName"],
 ["GET","user1","lastName"],
 ["GET","user1","city"],
 ["GET","lastName","user1"]]
*/
*/
