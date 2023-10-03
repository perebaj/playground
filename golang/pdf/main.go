package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/dslipak/pdf"
)

func main() {

	r, err := pdf.Open("teste.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println(r.NumPage())

	r2, err := r.GetPlainText()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("teste.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r2)

	f.WriteString(buf.String())
}

