/*
This code example is based on the talk of Rob Pike: https://youtu.be/f6kdp27TYZs?si=ztGWE5e-DszxMuru
All the credits to him
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Generator pattern: function that returns a channel
	joe := boring("Joe")
	ann := boring("Ann")
	for i := 0; i < 5; i++ {
		/* in this approach, we're stuck the ann because joe is blocking the channel(sequential execution), for this reason, we always
		will see the same output
		joe 1
		ann 1
		joe 2
		ann 2
		...
		To circumvent this problem, we can use the fan in pattern, where we can merge the channels into one channel
		c1
		\
			-----> c
		/
		c2

		In this way ann and joe will be completely independent
		*/
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}
	fmt.Println("You're boring; I'm leaving.")

	// Fan in

	f := fanIn(joe, ann)
	for i := 0; i < 10; i++ {
		fmt.Println(<-f)
	}
	fmt.Println("You're boring; I'm leaving.")

	// fan in with select statement
	f2 := fanIn2(joe, ann)
	for i := 0; i < 10; i++ {
		fmt.Println(<-f2)
	}
	fmt.Println("You're boring; I'm leaving.")

	/*
			// Timeout
			Here we're using the timeout pattern, where we're using the select statement with the time.After function
			In the end I'm saying that if the channel is blocked for more than 1 second, I'm leaving, this for each message

		b := boring("jojo")
		for {
			select {
			case s := <-b:
				fmt.Println(s)
			case <-time.After(1 * time.Second):
				fmt.Println("You're too slow")
				return
			}
		}

	*/

	// Timeout
	// But I'm also could say that the whole loop has a timeout, so here, I'm defining a timeout for the whole loop
	// b := boring("jojo")
	// timeout := time.After(2 * time.Second)
	// for {
	// 	select {
	// 	case s := <-b:
	// 		fmt.Println(s)
	// 	case <-timeout:
	// 		fmt.Println("You're too slow")
	// 		return
	// 	}
	// }

	// Quit channel
	// Here I'm using a quit channel to stop the boring function
	quit := make(chan bool)
	bq := boringquit("quit", quit)
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-bq)
	}
	quit <- true
}

// Generator pattern: function that returns a channel
func boring(msg string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

func fanIn2(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

func boringquit(msg string, quit chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case <-quit:
				return
			case c <- fmt.Sprintf("%s %d", msg, i):
				// do something
			}
		}
	}()
	return c
}
