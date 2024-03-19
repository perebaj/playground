package main

import "fmt"

func main() {
	strs := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	groupAnagrams(strs)
}

/*
	eat -> aet
	ate ->
*/

func groupAnagrams(strs []string) [][]string {
	group := make(map[Key][]string)
	for _, v := range strs {
		key := key(v)
		_, ok := group[key]
		if !ok {
			group[key] = []string{v}
		} else {
			group[key] = append(group[key], v)
		}
	}

	var resp [][]string
	for _, v := range group {
		resp = append(resp, v)
	}
	//nlog(n)
	fmt.Println(resp)
	return resp
}

type Key [26]byte

func key(str string) Key {
	var k Key

	for i := 0; i < len(str); i++ {
		k[str[i]-'a']++
	}
	return k
}
