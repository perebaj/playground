package main

import "fmt"

func zero(iptr *int) {
	// *iptr = 0
	*iptr = 0
}

func main() {
	i := 1
	zero(&i)
	fmt.Println(i)
}

func exercice2() {
	/*
		Voce recebe uma lista ordenada de forma crescente e deve encontrar qual numero está faltando

		Essa lista pode ter número positivos ou negativos

		input := []int{1, 2, 3, 5}                     //output: 4
		input2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} //output: 11
		intput3 := []int{-1, -2, -3, -4, -6}           //output: -5




		Essa pergunta eu não sabia ao certo, mas foi interessante pensar nisso

		Você tem uma função
		func zero(iptr *int) {
			*iptr = 0
		}

		func main() {
			i := 1
			zero(&i)
			fmt.Pritln(i) //0
		}

		Obvio que aqui o resultado é zero, mas a pergunta é sobre a atribuição de iptr na função zero

		Perguntas sobre DDL, DQL, DML, DCL e TCL
		Perguntas sobre solid SRP, OCP, LSP, ISP, DIP
	*/
}
