package main

import "fmt"

func main() {
	ransomNote := "bg"
	magazine := "efjbdfbdgfjhhaiigfhbaejahgfbbgbjagbddfgdiaigdadhcfcj"
	fmt.Println(canConstruct(ransomNote, magazine))
}

func canConstruct(ransomNote string, magazine string) bool {
	if len(ransomNote) > len(magazine) {
		return false
	}

	mMagazine := make(map[byte]int)
	ransomMap := make(map[byte]int)

	for i := 0; i < len(magazine); i++ {
		if i < len(ransomNote) {
			ransomMap[ransomNote[i]]++
		}
		mMagazine[magazine[i]]++
	}
	fmt.Println(mMagazine, ransomMap)
	for ramKey, ramVal := range ransomMap {
		_, ok := mMagazine[ramKey]
		if !ok || ramVal > mMagazine[ramKey] {
			return false
		}
	}

	return true
}
