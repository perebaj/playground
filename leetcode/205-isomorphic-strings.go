package main

import "fmt"

func main() {
	s := "foo"
	t := "bar"

	fmt.Println(isIsomorphic(s, t))
}

/*

	b: 4
	a: 3

	a: 4
	b: 3
*/

func isIsomorphic(s string, t string) bool {
	mS := make(map[byte]byte)

    if len(s)==1 && len(t)==1 {
        return true
    }

	// comparing s with t
	for i := 0; i < len(s); i++ {
		c := s[i]
		vv, sok := mS[c]
        fmt.Println(vv)
		if !sok && c != t[i] {
			mS[c] = t[i]
			fmt.Println("aklsdj")
		} else if vv == t[i] {
			continue
		} else {
			return false
		}
		// if !sok && vv != t[i] {
		// 	return false
		// } else {
		// 	mS[s[i]] = t[i]
		// }

	}

	fmt.Println(mS)

	return true
}
/*
	bar foo

	b -> f
	a -> o
	r -> o x

	for over len(t){
		vt, okt := tHM
		vs, oks := sHM

		if not okt and not oks {

		}
	}
*/
