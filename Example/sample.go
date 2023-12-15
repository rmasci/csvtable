package main

import (
	"fmt"
	"github.com/rmasci/csvtable"
	"os"
)

func main() {
	var render string
	if len(os.Args) <= 1 {
		render = "Render=mysql"
	} else {
		render = "Render=" + os.Args[1]
	}
	inFile, err := os.ReadFile("sample.csv")
	if err != nil {
		fmt.Println("ERR:", err)
		os.Exit(1)
	}
	out, err := csvtable.Table(string(inFile), render)
	if err != nil {
		fmt.Println("ERR:", err)
		os.Exit(1)
	}
	fmt.Println(out)
}
