package main

import "fmt"

func main() {
	ransomNote := "bg"
	magazine := "efjbdfbdgfjhhaiigfhbaejahgfbbgbjagbddfgdiaigdadhcfcj"

	fmt.Println(canConstruct(ransomNote, magazine))
}

func canConstruct(ransomNote string, magazine string) bool {
	magazineMap := make(map[rune]int)

	for _, v := range magazine {
		magazineMap[v]++
	}

	ransomMap := make(map[rune]int)
	for _, v := range ransomNote {
		ransomMap[v]++
	}

	fmt.Println(magazineMap)
	fmt.Println(ransomMap)

	if len(ransomMap) > len(magazine) {
		return false
	}

	for k := range ransomMap {
		_, ok := magazineMap[k]
		if !ok || magazineMap[k] < ransomMap[k] {
			return false
		}

	}

	return true
}
