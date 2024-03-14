package main

import "fmt"

func main() {
	fmt.Println(isValid("()"))      //true
	fmt.Println(isValid("()[]{}"))  //true
	fmt.Println(isValid("(]"))      //false
	fmt.Println(isValid("[]("))     //false
	fmt.Println(isValid("[]()}{}")) //false

}

func isValid(s string) bool {

	if !isOpen(string(s[0])) {
		return false
	}

	var stack []string
	for i := 0; i < len(s); i++ {
		if isOpen(string(s[i])) {
			stack = append(stack, string(s[i]))
		}
		if !isOpen(string(s[i])) && len(stack) != 0 {
			top := stack[len(stack)-1]
			if similar(top, string(s[i])) {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		} else if len(stack) == 0 && !isOpen(string(s[i])) {
			return false
		}
	}

	if len(stack) != 0 {
		return false
	}

	return true
}

func similar(open, close string) bool {
	if open == "{" && close == "}" {
		return true
	}
	if open == "[" && close == "]" {
		return true
	}
	if open == "(" && close == ")" {
		return true
	}
	return false
}

func isOpen(v string) bool {
	switch v {
	case "[", "{", "(":
		return true
	default:
		return false
	}
}
