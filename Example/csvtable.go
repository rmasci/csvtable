package main

import (
	"fmt"
	"os"

	csvtable "github.com/rmasci/csvformat"
)

func main() {
	var out string
	cBytes, err := os.ReadFile("cities.csv")
	if err != nil {
		fmt.Println("Error Reading File", err)
		os.Exit(1)
	}
	out, err = csvtable.Table(string(cBytes), "Render=simple")
	fmt.Println(out)
}
