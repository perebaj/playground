package main

import "fmt"

type S struct {
	PlayerID int    `json:"playerId"`
	Scout    string `json:"scout"`
}

func main() {
	in2 := []S{
		{
			PlayerID: 1,
			Scout:    "gol",
		},
		{
			PlayerID: 2,
			Scout:    "assistencia",
		},
		{
			PlayerID: 1,
			Scout:    "gol",
		},
	}

	entrada := []map[string]interface{}{
		{"playerId": 1, "scout": "gol"},
		{"playerId": 2, "scout": "assistencia"},
		{"playerId": 1, "scout": "gol"},
	}

	fmt.Println(entrada[0]["playerId"])

	fmt.Println(in2)
	m := make(map[int]map[string]int)

	m[1] = make(map[string]int)
	fmt.Println(m)

	// for _, v := range in2 {
	// 	_, ok2 := m[v.PlayerID]
	// 	if !ok2 {
	// 		m[v.PlayerID] = make(map[string]int)
	// 	}
	// 	fmt.Println(m)

	// 	_, ok := m[v.PlayerID][v.Scout]
	// 	if !ok {
	// 		m[v.PlayerID][v.Scout]++
	// 	}
	// }
}
