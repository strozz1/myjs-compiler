package st

import (
	"fmt"
	"io"
	"os"
)

type LexemeType string

// Defines all the possible types of a lexeme.
// The 'LEXEME NONE' represents a lexeme that either hasn't been asigned yet
// or that it doesn't have a type.
const (
	LEXEME_FUNCTION  = "funcion"
	LEXEME_PROCEDURE = "procedimiento"
	LEXEME_INTEGER   = "entero"
	LEXEME_STRING    = "cadena"
	LEXEME_REAL      = "real"
	LEXEME_LOGIC     = "logico"
	LEXEME_POINTER   = "puntero"
	LEXEME_VECTOR    = "vector"
	LEXEME_NONE      = "none"
)

type Entry struct {
	Lexeme     string                //lexeme
	Type       LexemeType            //lexeme type
	Id         int                   //id of Entry
	Attributes map[string]*Attribute //attribute list
}

// Creates new symbol entry from the original lexeme.
// Returns the new Entry object with attributes init.
//
// Sets the type to LEXEME NONE by default
func NewEntry(lex string) *Entry {
	return &Entry{
		Attributes: map[string]*Attribute{},
		Lexeme:     lex,
		Type:       LEXEME_NONE,
	}

}

// sets the type of the lexeme
// IF an invalid lexem type is provided, an error is returned
func (e *Entry) SetType(t LexemeType) error {
	switch t {
	case LEXEME_FUNCTION:
		break
	case LEXEME_PROCEDURE:
		break
	case LEXEME_INTEGER:
		break
	case LEXEME_STRING:
		break
	case LEXEME_REAL:
		break
	case LEXEME_LOGIC:
		break
	case LEXEME_POINTER:
		break
	case LEXEME_VECTOR:
		break
	case LEXEME_NONE:
		break
	default:
		{
			if DEBUG {
				fmt.Printf("DEBUG: Invalid Lexeme Type: [%v]\n\r", t)
			}
			return fmt.Errorf("Error: Invalid Lexem type: [%v]", t)
		}
	}
	e.Type = t
	return nil
}

// Writes the SymbolEntry from the ST to the specified Writer with
// PDL specified format
func (e *Entry) Write(w io.Writer) {
	fmt.Fprintf(w, "* LEXEMA: '%v'\r\n", e.Lexeme)
	fmt.Fprintln(w, "  Atributos:")
	fmt.Fprintf(w, "    + Tipo: ")
	if e.Type == LEXEME_NONE {
		fmt.Fprintf(w, "'-'")
	} else {
		fmt.Fprintf(w, "'%v'", e.Type)
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
			attribute.StringVal = fmt.Sprintf("%v", value)
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
			attribute.IntVal = v
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
			attribute.ArrayVal = v
		}
	default:
		if DEBUG {
			fmt.Printf("DEBUG: Invalid Attribute type: [%v]\n\r", attribute.Type)
			return
		}
	}
	attribute.HasValue = true
}

// Adds an attribute to the Entry.
// Returns error if the attribute is already present
func (e *Entry) AddAtribute(name string, a Attribute) error {
	if e.containsAttr(name) {
		return fmt.Errorf("Error: Attribute '%v' already exists for the entry [%v]", name, e.Lexeme)
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
