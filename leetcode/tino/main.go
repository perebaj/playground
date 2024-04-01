/*
Você está escrevendo uma parte do nosso sistema anti-fraude de pagamentos.

Uma das formas de detectarmos fraudes é quando o mesmo cliente faz mais de um pagamento numa janela muito curta de tempo.

Seu programa receberá os clientes que estão efetuando pagamentos através de um stream, e deve identificar clientes repetidos que apareçam próximos.

Escreva um algoritmo que leia um stream de clientes, onde cada cliente é identificado por um número inteiro, e escreva na saída cada vez que um cliente aparecer repetido numa janela de tamanho J.

**Exemplo:**
J = 4

Stream = [7, 10, 5, 10, 8, 3, 1, 4, 3, 3, 5, 1]

Saída: 10, 3, 3
*/

package main

import "fmt"

func main() {
	var stack []int // len == j remove fist element
	stream := []int{7, 10, 5, 10, 8, 3, 1, 4, 3, 3, 5, 1} //n = numeros de elementos do stream
	j := 4
	for _, v := range stream {
		stack = append(stack, v)
		if len(stack) == j {
			stack = stack[1:]
			m := make(map[int]int)
			for _, v := range stack {
				m[v]++
			}

			for k, v := range m {
				if v > 1 {
					fmt.Println(k)
				}
			}
		}
	}
}
