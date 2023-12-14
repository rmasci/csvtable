package csvtable

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/rmasci/gotabulate"
	"strconv"
	"strings"
	"time"
)

// Grid Render Used for Gridout function.  Renders output text in a grid, text, tab.
//
//		Render:
//			simple Simple tab delimited with underlined header.
//			plain pretty much just strips the "," for a space
//			tab   Tab spaced output
//			html  Output in simple HTML. Header has grey background with thin lines between cells
//			mysql Looks similar to mysql shell output
//			grid  uses ANSI Graphics. Not compatible with all terminals
//			gridt MySQL
//		Align: Right, left or center
//		NoHeader: don't print a header
//		Wrap: wrap cell output
//		NoLineBetweenRow: put a blank line inbetween each row
//		Columns: List the columns you want output.
//		Space: Do not trim spaces from column. "1,   2" will output "1|    2" instead of "1|2"
//		Delimeter: Delimeter between text. Default is a ","
//	 Number -- automatically add number the output
type Grid struct {
	Render           string
	Align            string
	NoHeader         bool
	Wrap             bool
	Indent           string
	NoLineBetweenRow bool
	Space            bool
	Columns          string
	Delimiter        string
	Number           bool
	OutDelimiter     string
	Headline         string
}

// TimeStamp
// Returns the current time as timestamp: %Y%m%d%H%M%s  optionally you can give it a time format in accordance with Go's time package.
func TimeStamp(a ...any) string {
	var f string
	if len(a) >= 1 {
		f = fmt.Sprintf("%v", a[0])
	} else {
		f = "20060102150405"
	}
	return time.Now().Format(f)
}

// NewGrid Returns a Grid with some defaults.
//
//	Render: mysql
//	Align: left
//	NoHeader: False
//	Wrap: false
//	Indent: "" (empty)
//	LineBetweenRow: false
//	Columns: All (empty)
//
// New grid gives defaults. You can change these after NewGrid or just set it using the struct as shown in NewGrid.
// your own use might look like this
// mygrid := nsres.Grid{Render: "grid", Align: "Right", NoHeader: false, Wrap: false, Indent: "\t", LineBetweenRow: false, Space: false, Columns: "", Delimiter: "|"
func NewGrid() *Grid {
	g := Grid{Render: "mysql", Align: "Left", NoHeader: false, Wrap: false, Indent: "", NoLineBetweenRow: true, Space: true, Columns: "", Delimiter: ","}
	return &g
}

// Table creates a formatted table from the provided text string.
// The function accepts an optional list of flags that can be used to customize the table's formatting.
// The flags can be used to specify the table's render style, alignment, header presence, wrapping, indentation, line separation between rows, spacing, columns, and delimiter.
// If the flags are not provided, the table is created with default settings.
// The function returns the formatted table as a string and an error if any occurred during the table creation.
// Try these with <render>-nohead when there is no header.
//
//	+-----------+-------------------------------------+
//	| Render    | Output Format                       |
//	+===========+=====================================+
//	| mysql     | Looks like a MySQL Client Query     |
//	| grid      | Spreadsheet using Graphical Grid    |
//	| gridt     | Spreadsheet using text grid         |
//	| simple    | Simple Table                        |
//	| html      | Output in HTML Table                |
//	| tab       | Just text tab separated             |
//	| csv       | Output in CSV format                |
//	| plain     | Plain Table output                  |
//	+-----------+-------------------------------------+
func Table(text string, flags ...string) (string, error) {
	g := NewGrid()
	if len(flags) != 0 {
		err := g.parseFlags(flags...)
		if err != nil {
			return "", err
		}
	}
	switch g.Render {
	case "grid":
		fallthrough
	case "grid-nohead":
		fallthrough
	case "mysqlg-nohead":
		fallthrough
	case "mysqlg":
		g.NoLineBetweenRow = false
	}
	return g.Gridout(text)
}

func (g *Grid) parseFlags(flags ...string) error {
	if len(flags) == 0 {
		return nil
	}
	var err error
	for _, f := range flags {
		if !strings.Contains(f, "=") {
			continue
		}
		arg := strings.Split(f, "=")

		switch arg[0] {
		case "Render":
			g.Render = arg[1]
		case "Align":
			g.Align = arg[1]
		case "NoHeader":
			g.NoHeader, err = stringToBool(arg[1])
			if err != nil {
				return err
			}
		case "Wrap":
			g.Wrap, err = stringToBool(arg[1])
			if err != nil {
				return err
			}
		case "NoLineBetweenRow":
			g.NoLineBetweenRow, err = stringToBool(arg[1])
			if err != nil {
				return err
			}

		case "Space":
			g.Space, err = stringToBool(arg[1])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Gridout
// Don't waste time getting the output looking nice, just make it csv, then pass it to gridoutl.  Gridout Prints according to Grid struct.
// Render:
// simple. Tab format with lines on top
// plain. Text output no gridlines extra spaces in between
func (g *Grid) Gridout(text string) (string, error) {
	var out string
	//if g.Render == "text" {
	//	//g.Render = "csv"
	//	g.Render = "txt"
	//}
	if g.Headline != "" {
		strings.TrimSpace(g.Headline)
		text = fmt.Sprintf("%s\n%s", g.Headline, text)
	}
	txtRdr := strings.NewReader(text)
	csvRdr := csv.NewReader(txtRdr)
	csvRdr.Comma = []rune(g.Delimiter)[0]
	csvRdr.TrimLeadingSpace = g.Space
	s, err := csvRdr.ReadAll()
	if len(s) <= 1 {
		g.NoHeader = true
	}
	if err != nil {
		return "", fmt.Errorf("Error csvRdr.ReadAll: %v\n", err)
	}
	// Normalize for errors since we get an error trying to convert a non number to string.
	g.Columns = strings.TrimSpace(g.Columns)
	for i, _ := range s {
		if g.Columns != "" {
			var tmp []string
			pCol := strings.Split(g.Columns, ",")
			for _, c := range pCol {
				c = strings.TrimSpace(c)
				idx, err := strconv.Atoi(c)
				if err != nil {
					//return "", fmt.Errorf("can't convert %v to int %err", c, err)
					continue
				}
				if idx > len(s[i]) || idx <= 0 {
					err := fmt.Errorf("invalid column col:%v len:%v", idx, len(s[i]))
					return "", err
				}
				idx = idx - 1
				if g.Render == "bingo" {
					s[i][idx] = fmt.Sprintf("\n\n%s\n\n", s[i][idx])
				}
				tmp = append(tmp, s[i][idx])
			}
			s[i] = tmp
		}
		if i == 0 && g.Number {
			s[i] = append([]string{"Index"}, s[i]...)
		} else if g.Number {
			I := fmt.Sprintf("%d", i)
			s[i] = append([]string{I}, s[i]...)
		}
	}
	gridulate := gotabulate.Create(s)
	gridulate.SetAlign(g.Align)
	gridulate.SetWrapStrings(g.Wrap)
	gridulate.SetRemEmptyLines(g.NoLineBetweenRow)
	gridulate.NoHeader = g.NoHeader
	scanner := bufio.NewScanner(strings.NewReader(gridulate.Render(g.Render)))
	for scanner.Scan() {
		out = fmt.Sprintf("%s%s%s\n", out, g.Indent, scanner.Text())
	}
	return out, nil
}

func stringToBool(s string) (bool, error) {
	s = strings.ToLower(s)
	boolValue, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return boolValue, nil
}
