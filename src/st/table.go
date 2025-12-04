package st

import (
	"fmt"
	"io"
)

type SymbolTable struct {
	id     int
	name   string
	table  map[string]*Entry
	keys   []string
	inner  *SymbolTable
	parent *SymbolTable
	offset int
}

// Interal function to create a new SymbolTable.
// Initializes the table and assings an id
func createST(name string) *SymbolTable {
	stIdCounter++
	return &SymbolTable{
		id:    stIdCounter,
		name:  name,
		keys:  []string{},
		offset: 0,
		table: map[string]*Entry{},
	}

}

func (s *SymbolTable) GetEntry(pos int) (*Entry, bool) {
	if len(s.keys) <= pos {
		return nil, false
	}
	name := s.keys[pos]
	a, ok := s.table[name]
	return a, ok
}

// Adds a new Symbol/Entry to the table. If it already exists returns a nil.
func (s *SymbolTable) AddEntry(lex string) (int, bool) {
	e, err := s.table[lex]
	if err {
		if DEBUG {
			fmt.Printf("DEBUG: Failed to insert already existing Symbol '%v' on table [%v]\n\r", lex, s.name)
		}
		return -1, false
	}
	e = NewEntry(lex)
	s.table[lex] = e
	s.keys = append(s.keys, lex)
	_, ok := s.table[lex]
	if DEBUG {
		fmt.Printf("DEBUG: Added new entry '%v' to table '%v'\n\r", lex, s.name)
		if !ok {
			fmt.Printf("ERROR: cant add entry '%v' to table %s\n", lex, s.name)
		}
	}
	e.pos = len(s.keys) - 1
	return e.pos, ok
}

func (s *SymbolTable) RemoveEntry(lex string) {
	delete(s.table, lex)
	if DEBUG {
		fmt.Printf("DEBUG: Removed entry '%v' from table '%v'\n\r", lex, s.name)
	}
}

// Writes the SymbolTable in the Specified format for PDL.
// @input: io.Writer
func (s *SymbolTable) Write(w io.Writer) {
	if DEBUG {
		fmt.Printf("DEBUG: Writing table '%v' to output\n\r", s.name)
	}
	fmt.Fprintf(w, "%v #%d:\n\r", s.name, s.id)
	for _, i := range s.table {
		i.Write(w)
	}
	fmt.Fprintln(w, "------------------------------------------")
}
