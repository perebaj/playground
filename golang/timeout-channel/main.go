package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	timeout := time.After(1 * time.Second)
	c := workInProgress("hello")

	for {
		select {
		case <-timeout:
			println("timeout")
			os.Exit(0)
		case s := <-c:
			println(s)
		}
	}
}

func workInProgress(msg string) chan string {
	c := make(chan string)

	go func() {
		for {
			c <- fmt.Sprintf("work in progress: %s", msg)
		}
	}()

	return c
}
