package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dirfs := home + "/playground/golang/fs" // absolute path instead of relative to run this program from anywhere
	dir := os.DirFS(dirfs)
	info, err := fs.Stat(dir, ".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Name:", info.Name())
}
