package main

import (
	"fmt"
	"unicode"
)

/*
65 - 90  UPPER CASE (convert)
97 - 122  LOWER CASE
all other ignore
*/
func main() {
	// fmt.Println(byte["a"])
	s := "A man, a plan, a s canal: Panama"
	fmt.Println(isPalindrome(s))
	// fmt.Println(convert(s))

}

func isPalindrome(s string) bool {
	convertedString := convert(s)
	initial := 0
	end := len(convertedString) - 1

	for initial < end {
		if convertedString[initial] != convertedString[end] {
			return false
		}
		initial++
		end--
	}
	return true
}

func convert(s string) string {
	var resp string
	for _, v := range s {
		if unicode.IsLetter(v) || unicode.IsNumber(v) {
			if unicode.IsUpper(v) {
				resp = resp + string((v + 32))
			} else {
				resp = resp + string(v)
			}
		}
	}

	return resp
}

/*
Here I learned something new about non-alpanumeric and alphanumir, it's essencial to solve it!
and also, the unicode library

alphanumeric contempalte all letters and digits

*/
