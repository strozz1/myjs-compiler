package st

import (
	"fmt"
	"io"
	"os"
)

type LexemeKind int

// Defines all the possible types of a lexeme.
// The 'LEXEME NONE' represents a lexeme that either hasn't been asigned yet
// or that it doesn't have a type.
const (
	LEXEME_NONE LexemeKind = iota
	LEXEME_FUNCTION
	LEXEME_PROCEDURE
	LEXEME_INTEGER
	LEXEME_STRING
	LEXEME_REAL
	LEXEME_LOGIC
	LEXEME_POINTER
	LEXEME_VECTOR
)

func (l LexemeKind) String() string {
	var str string
	switch l {
	case LEXEME_FUNCTION:
		str = "funcion"
	case LEXEME_PROCEDURE:
		str = "procedimiento"
	case LEXEME_INTEGER:
		str = "entero"
	case LEXEME_STRING:
		str = "cadena"
	case LEXEME_REAL:
		str = "real"
	case LEXEME_LOGIC:
		str = "logico"
	case LEXEME_POINTER:
		str = "puntero"
	case LEXEME_VECTOR:
		str = "vector"
	default:
		str = "indefinido"
	}
	return str
}

type Entry struct {
	id         int                   //id of Entry
	lexeme     string                //lexeme
	kind       LexemeKind            //lexeme kind
	Attributes map[string]*Attribute //attribute list
	Pos        int                   // pos
}

// Creates new symbol entry from the original lexeme.
// Returns the new Entry object with attributes init.
//
// Sets the type to LEXEME NONE by default
func NewEntry(lex string) *Entry {
	return &Entry{
		Attributes: map[string]*Attribute{},
		lexeme:     lex,
		kind:       LEXEME_NONE,
	}

}

// sets the type of the lexeme
// IF an invalid lexem type is provided, an error is returned
func (e *Entry) SetType(t LexemeKind) error {
	switch t {
	case LEXEME_FUNCTION:
	case LEXEME_PROCEDURE:
	case LEXEME_INTEGER:
	case LEXEME_STRING:
	case LEXEME_REAL:
	case LEXEME_LOGIC:
	case LEXEME_POINTER:
	case LEXEME_VECTOR:
	case LEXEME_NONE:
	default:
		{
			if DEBUG {
				fmt.Printf("DEBUG: Invalid Lexeme Type: [%v]\n\r", t)
			}
			return fmt.Errorf("Error: Invalid Lexem type: [%v]", t)
		}
	}
	e.kind = t
	return nil
}

// Writes the SymbolEntry from the ST to the specified Writer with
// PDL specified format
func (e *Entry) Write(w io.Writer) {
	fmt.Fprintf(w, "* LEXEMA: '%v'\r\n", e.lexeme)
	fmt.Fprintln(w, "  Atributos:")
	fmt.Fprintf(w, "    + Tipo: ")
	if e.kind == LEXEME_NONE {
		fmt.Fprintf(w, "'-'")
	} else {
		fmt.Fprintf(w, "'%v'", e.kind.String())
	}
	fmt.Fprintln(w)
	for _, at := range e.Attributes {
		at.Write(w)
	}

}

// Adds a value to the attribute specified with param 'name'.
// If it succeeds, it sets 'HasValue' of the attribute to true.
// The attribute needs to exist before asigning a value.
func (e *Entry) SetAttributeValue(name string, value any) {
	if !e.containsAttr(name) {
		if DEBUG {
			fmt.Printf("DEBUG: Can't add value to a non existing attribute [%v]\n\r", name)
			return
		}
	}
	attribute, _ := e.Attributes[name]
	switch attribute.Type {
	case T_STRING:
		{
			attribute.stringVal = fmt.Sprintf("%v", value)
			break
		}
	case T_INTEGER:
		{
			v, ok := value.(int)
			if !ok {
				if DEBUG {
					fmt.Printf("DEBUG: Invalid value for type Integer: [%v]\n\r", value)
					return
				}
			}
			attribute.intVal = v
			break
		}
	case T_ARRAY:
		{
			v, ok := value.([]string)
			if !ok {
				if DEBUG {
					fmt.Printf("DEBUG: Invalid value for type array: [%v]\n\r", value)
					return
				}
			}
			attribute.arrayVal = v
		}
	default:
		if DEBUG {
			fmt.Printf("DEBUG: Invalid Attribute type: [%v]\n\r", attribute.Type)
			return
		}
	}
	attribute.hasValue = true
}

// Adds an attribute to the Entry.
// Returns error if the attribute is already present
func (e *Entry) AddAtribute(name string, a Attribute) error {
	if e.containsAttr(name) {
		return fmt.Errorf("Error: Attribute '%v' already exists for the entry [%v]", name, e.lexeme)
	}
	e.Attributes[name] = &a
	return nil
}

// Returs a bool if val is in the attributes
func (e *Entry) containsAttr(val string) bool {
	_, v := e.Attributes[val]
	return v
}

func printTable(table map[string]*Attribute) {
	if DEBUG {
		fmt.Println("Printing table:")
		for _, k := range table {
			k.Write(os.Stdout)
		}
	}
}
