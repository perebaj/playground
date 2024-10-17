package main

import (
	"fmt"w
	"time"
)

func main() {
	jjTest("jojo")
	time.Sleep(time.Second * 1)
}

func jjTest(s string) {
	if s == "jojo" {
		fmt.Println("jojo without goroutine")
	} else {
		go func() {
			fmt.Println("jojo with goroutine")
		}()
	}

}
