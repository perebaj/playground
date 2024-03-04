package main

import "fmt"

func main() {
	block := "WBBWWBBWBW"
	k := 7
	fmt.Println(minimumRecolors(block, k))
}

func minimumRecolors(blocks string, k int) int {
	result := 101
	var aux int
	for i := 0; i < len(blocks)-k+1; i++ {
		sliceAux := blocks[i : k+i]
		for _, v := range sliceAux {
			if string(v) == "W" {
				aux++
			}
		}
		if aux < result {
			result = aux
		}
		fmt.Println(aux)
		aux = 0
	}
	return result
}

/*
   blocks = "WBBWWBBWBW", k = 7
   W B B W W B B W B W

   1) W B B W W B B
   2) B B W W B B W
   3) B W W B B W B
   4) W W B B W B W

   LEN(BLOCK) - K + 1

   if number_of_white == 0 {
       return 0
   } else if aux_number_of_white < number_of_white  {
       number_of_white = aux_num_of_white
   }

   I think to use the SQ approach, we could calculate the

   1) calculate the first substring and save if
   2) after that we can just calculate the corners and compare with the current result
     B W
	 B B -1

	 B w
	 w w +1

	 B B
	 B W +1

	 W W
	 B B -2

	 B B
	 W W
*/
