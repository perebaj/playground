package main

func main() {
	canCompleteCircuit(nil, nil)
}

func canCompleteCircuit(gas []int, cost []int) int {
	totalSurplus, currentSurplus := 0, 0
	startStation := 0
	for i := 0; i < len(gas); i++ {
		totalSurplus = totalSurplus + gas[i] - cost[i]
		currentSurplus = currentSurplus + gas[i] - cost[i]
		if currentSurplus < 0 {
			startStation = i + 1
			currentSurplus = 0
		}
	}

	// If the totoal surplus is greater than 0
	// there's a way to travel around the circuit
	if totalSurplus >= 0 {
		return startStation
	} else {
		return -1
	}
}
/*
Eu não consegui ter a sacada p/ resolver esse exercicio, de fato, algo bem tricky

Mas eu entendi a solução no final

ALgo do tipo:
a diferença inicial entre, gas e custo na posição i, não pode ser menor que zero, se for, você basicamente está começando com combustivel negativo, o
que não faz o menor sentido

portanto, temos q encontrar algum ponto onde isso não seja verdade!

O que eu estava travado, é se eu precisaria verificar isso circularmente, mas a sacada é essa, não!

se você validar que a sum(gas) == sum(cost), então eu só preciso validar uma posição que a diferença de gas-cost n é negativa, e validar os valores subsequentes


dado o exemplo
[1,2,3,4,5]
[3,4,5,1,2]
[-2,-2,-2,3,3]

a unica posição valida para começar é a 3

ou seja nossa tanque tem 3 de combustivel

tank = 3

para ir para o proximo passo, gastamos 2, ou seja 3-2 = 1 e ao chegar nessa estação, carregamos com 5, ou seja 8

8 > 0, sucesso!
**/
