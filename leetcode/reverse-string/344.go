package main

func main() {
	input := []byte{'h', 'e', 'l', 'l', 'o'}
	reverseString(input)
}

func reverseString(s []byte) {
	startPointer := 0
	endPointer := len(s) - 1
	for startPointer < endPointer {
		startV := s[startPointer]
		endV := s[endPointer]
		s[startPointer] = endV
		s[endPointer] = startV
		startPointer++
		endPointer--
	}
}
