package main

import (
	"fmt"
	"strconv"
)

func main() {
	tokens := []string{"3", "11", "+", "5", "-"} // -9 9
	// fmt.Println(evalRPN(tokens))
	fmt.Println(evalRPN2(tokens))
}

// There is a bug when using the - operand that I don't know how to solve
// I perceive that to solve it, we don't need to create 2 stack, only one it's enogh
func EvalRPN(tokens []string) int {
	if len(tokens) == 1 {
		integer, _ := strconv.Atoi(tokens[0])
		return integer
	}
	var resultStack []int
	var auxStack []int
	var result int
	for i := 0; i < len(tokens); i++ {

		ok := isOperand(tokens[i])
		if !ok {
			integer, _ := strconv.Atoi(tokens[i])
			auxStack = append(auxStack, integer)
		} else if ok && len(resultStack) == 0 {
			v2 := auxStack[len(auxStack)-1]
			//pop
			auxStack = auxStack[:len(auxStack)-1]
			v1 := auxStack[len(auxStack)-1]
			//pop
			auxStack = auxStack[:len(auxStack)-1]
			result = operation(v1, v2, tokens[i])
			resultStack = append(resultStack, result)
		} else {
			v1 := resultStack[len(resultStack)-1]
			resultStack = resultStack[:len(resultStack)-1]
			v2 := auxStack[len(auxStack)-1]
			auxStack = auxStack[:len(auxStack)-1]
			result = operation(v1, v2, tokens[i])
			resultStack = append(resultStack, result)
		}
	}

	fmt.Println(resultStack, auxStack, result)
	return result
}

func isOperand(b string) bool {
	switch b {
	case "+", "/", "*", "-":
		return true
	default:
		return false
	}
}

func operation(v1, v2 int, operan string) int {
	switch operan {
	case "+":
		return v1 + v2
	case "-":
		return v1 - v2
	case "/":
		return v1 / v2
	case "*":
		return v1 * v2
	}
	return 0
}

// THis one was approved by all tests
func evalRPN2(tokens []string) int {
	var stack []int
	for i := 0; i < len(tokens); i++ {
		integer, err := strconv.Atoi(tokens[i])
		if err == nil {
			// append
			stack = append(stack, integer)
		} else {
			value2 := stack[len(stack)-1]
			value1 := stack[len(stack)-2]
			result := operation(value1, value2, tokens[i])
			stack = stack[:len(stack)-2]
			stack = append(stack, result)
		}
		fmt.Println(stack)
	}
	return 0
}
