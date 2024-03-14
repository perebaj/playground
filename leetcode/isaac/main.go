/*
Dada uma lista de número inteiros `nums` e um número inteiro `target`, retorne os índices de dois números cuja soma é igual ao `target`.

Você pode assumir que cada entrada teria *exatamente uma solução*, e você não pode usar o mesmo elemento duas vezes.

Você pode retornar a resposta em qualquer ordem.
//Input:
nums = [2,7,11,15], target = 9

//Output:
[0,1]

//Input:
nums = [3,2,4], target = 6

// Output:
Output: [1,2]

 1. fixar um ponto no array

 2. somar com todos os outros elementos
    2.1) validar se essa soma == target
    2.1) se for retorna indice do ponteiro fixado + indice que esta correndo

    0 x 1
    0 x 2
    0 x 3
    1 x 2
    1 x 3
    ...
*/
package main

import "fmt"

func main() {
	input := []int{3, 3}
	target := 6
	fmt.Println(resolution2(input, target))
}

func resolution(s []int, t int) []int {
	resp := make([]int, 2)
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i]+s[j] == t {
				fmt.Println(i, j)
				resp[0] = i
				resp[1] = j
			}
		}
	}
	return resp
}

func resolution2(s []int, t int) []int {
	m := make(map[int]int)
	for i, v := range s {
		m[v] = i
	}
	resp := make([]int, 2)

	for i := 0; i < len(s); i++ {
		aux := t - s[i]
		_, ok := m[aux]
		if ok {
			fmt.Print(i, m[aux])
			resp[0] = i
			resp[1] = m[aux]
			break
		}
	}
	return resp
}
