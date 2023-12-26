package main

import (
	"fmt"

	"github.com/gosimple/slug"
)

func main() {
	txt := slug.Make("i like you a lot")
	fmt.Println(txt)
}
