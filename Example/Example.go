package main

import (
	"fmt"
	"github.com/rmasci/csvtable"
	"os"
)

func main() {
	// Notice: No header. You can specify that with 'Header'
	csvtext := fmt.Sprintf("One,Two,Three\nFour,Five,Six\nSeven,Eight,Nine")
	csvgrid := csvtable.NewGrid()
	csvgrid.Header = "First,Second,Third" // This is optional. If the CSV already contains a heading line. don't set this.
	if len(os.Args) <= 1 {
		csvgrid.Render = "simple"
	} else {
		csvgrid.Render = os.Args[1]
	}
	out, err := csvgrid.Gridout(csvtext)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Println(out)
}
