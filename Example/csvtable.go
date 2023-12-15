package main

import (
	"fmt"
	"os"

	"github.com/rmasci/csvtable"
)

func main() {
	var out string
	cBytes, err := os.ReadFile("cities.csv")
	if err != nil {
		fmt.Println("Error Reading File", err)
		os.Exit(1)
	}
	g := csvtable.NewGrid()
	if len(os.Args) <= 1 {
		g.Render = "mysql"
	} else {
		g.Render = os.Args[1]
	}
	out, err = g.Gridout(string(cBytes))
	fmt.Println(out)
}
