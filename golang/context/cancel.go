package main

import (
	"context"
	"fmt"
	"time"
)

func working(ctx context.Context) error {
	count := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			count++
			fmt.Println(count)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err := working(ctx)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
