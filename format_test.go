package csvtable

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTimeStamp(t *testing.T) {
	// Test default format
	result := TimeStamp()
	if len(result) != 14 {
		t.Errorf("TimeStamp() with no arguments returned a string of length %d, want %d", len(result), 14)
	}

	// Test custom format
	customFormat := "2006-01-02"
	result = TimeStamp(customFormat)
	expected := time.Now().Format(customFormat)
	if result != expected {
		t.Errorf("TimeStamp() with custom format returned %s, want %s", result, expected)
	}
}

func TestNewGrid(t *testing.T) {
	grid := NewGrid()

	if grid.Render != "mysql" {
		t.Errorf("NewGrid() returned grid with Render = %s, want mysql", grid.Render)
	}

	if grid.Align != "Left" {
		t.Errorf("NewGrid() returned grid with Align = %s, want Left", grid.Align)
	}

	if grid.NoHeader != false {
		t.Errorf("NewGrid() returned grid with NoHeader = %v, want false", grid.NoHeader)
	}

	if grid.Wrap != false {
		t.Errorf("NewGrid() returned grid with Wrap = %v, want false", grid.Wrap)
	}

	if grid.Indent != "" {
		t.Errorf("NewGrid() returned grid with Indent = %s, want \"\"", grid.Indent)
	}

	if grid.NoLineBetweenRow != true {
		t.Errorf("NewGrid() returned grid with NoLineBetweenRow = %v, want true", grid.NoLineBetweenRow)
	}

	if grid.Space != true {
		t.Errorf("NewGrid() returned grid with Space = %v, want true", grid.Space)
	}

	if grid.Columns != "" {
		t.Errorf("NewGrid() returned grid with Columns = %s, want \"\"", grid.Columns)
	}

	if grid.Delimiter != "," {
		t.Errorf("NewGrid() returned grid with Delimiter = %s, want ,", grid.Delimiter)
	}
}

func TestParseFlags(t *testing.T) {
	g := NewGrid()

	err := g.parseFlags("Render=grid", "Align=Right", "NoHeader=true", "Wrap=true", "NoLineBetweenRow=f", "Space=true", "Number=true")
	if err != nil {
		t.Errorf("parseFlags() returned error: %v", err)
	}

	if g.Render != "grid" {
		t.Errorf("parseFlags() did not correctly set Render, got: %s, want: grid", g.Render)
	}

	if g.Align != "Right" {
		t.Errorf("parseFlags() did not correctly set Align, got: %s, want: Right", g.Align)
	}
	if g.NoHeader != true {
		t.Errorf("parseFlags() did not correctly set NoHeader, got %v, want True", g.NoHeader)
	}
	if g.Wrap != true {
		t.Errorf("parseFlags() did not correctly set Wrap, got: %v, want: true", g.Wrap)
	}

	if g.NoLineBetweenRow != false {
		t.Errorf("parseFlags() did not correctly set NoLineBetweenRow, got: %v, want: false", g.NoLineBetweenRow)
	}

	if g.Space != true {
		t.Errorf("parseFlags() did not correctly set Space, got: %v, want: true", g.Space)
	}
}

func TestTable(t *testing.T) {
	data := `Name,Age
Alice,30
Bob,25`
	table, err := Table(data)
	if err != nil {
		t.Errorf("Table() returned error: %v", err)
	}

	if !strings.Contains(table, "Name") || !strings.Contains(table, "Age") {
		t.Errorf("Table() did not correctly create table, headers are missing")
	}

	if !strings.Contains(table, "Alice") || !strings.Contains(table, "30") {
		t.Errorf("Table() did not correctly create table, Alice's data is missing")
	}

	if !strings.Contains(table, "Bob") || !strings.Contains(table, "25") {
		t.Errorf("Table() did not correctly create table, Bob's data is missing")
	}
}

func TestGridOut(t *testing.T) {
	g := NewGrid()
	data := `Name,Age
	Alice,30
	Bob,25`

	out, err := g.Gridout(data)
	if err != nil {
		t.Errorf("Gridout() returned error: %v", err)
	}

	if !strings.Contains(out, "Name") || !strings.Contains(out, "Age") {
		t.Errorf("Gridout() did not correctly create output, headers are missing")
	}

	if !strings.Contains(out, "Alice") || !strings.Contains(out, "30") {
		t.Errorf("Gridout() did not correctly create output, Alice's data is missing")
	}

	if !strings.Contains(out, "Bob") || !strings.Contains(out, "25") {
		t.Errorf("Gridout() did not correctly create output, Bob's data is missing")
	}
}

func TestStringToBool(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
		err      error
	}{
		{"true", true, nil},
		{"True", true, nil},
		{"TRUE", true, nil},
		{"false", false, nil},
		{"False", false, nil},
		{"FALSE", false, nil},
		{"invalid", false, fmt.Errorf(`strconv.ParseBool: parsing "invalid": invalid syntax`)},
		{"", false, fmt.Errorf(`strconv.ParseBool: parsing "": invalid syntax`)},
	}

	for _, tc := range testCases {
		result, err := stringToBool(tc.input)
		if err != nil && err.Error() != tc.err.Error() {
			t.Errorf("stringToBool(%q) returned error: %v, want: %v", tc.input, err, tc.err)
		}
		if result != tc.expected {
			t.Errorf("stringToBool(%q) = %v; want %v", tc.input, result, tc.expected)
		}
	}
}
