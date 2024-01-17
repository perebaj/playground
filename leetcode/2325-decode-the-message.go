package main

import "fmt"

func main() {
	key := "the quick brown fox jumps over the lazy dog"
	message := "vkbs bs t suepuv"
	fmt.Println(decodeMessage(key, message))
}

func decodeMessage(key string, message string) string {
	m := make(map[string]string)
	aux := 97

	for i := 0; i < len(key); i++ {
		if string(key[i]) == " " {
			continue
		}
		_, ok := m[string(key[i])]
		if !ok {
			m[string(key[i])] = string(aux)
			aux++
		}
	}

	var res string

	for i := 0; i < len(message); i++ {

		res = res + m[string(message[i])]
		if string(message[i]) == " " {
			res = res + " "
		}
	}

	fmt.Println(res)

	return res
}
