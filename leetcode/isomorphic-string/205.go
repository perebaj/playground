package main

import "fmt"

func main() {
	s := "badc"
	t := "baba"

	fmt.Println(isIsomorphic(s, t))
}

func isIsomorphic(s string, t string) bool {
	mS := make(map[string]string)
	mT := make(map[string]string)

	if len(mS) != len(mT) {
		return false
	}

	for i := 0; i < len(s); i++ {
		vS, ok := mS[string(s[i])]
		_, ok2 := mT[string(t[i])]
		fmt.Println(ok, mS, mT, string(s[i]), vS)
		fmt.Println(mT[vS], string(s[i]))
		fmt.Println(mT[vS], vS)
		if (ok2 || ok) && mT[string(t[i])] != string(s[i]) && mS[string(s[i])] != string(t[i]) {
			return false
		} else {
			mS[string(s[i])] = string(t[i])
			mT[string(t[i])] = string(s[i])
		}
	}
	return true
}
