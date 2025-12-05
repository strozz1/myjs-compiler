package st

import (
	"fmt"
)

type EntryType int

// Defines all the possible types of a lexeme.
// The 'LEXEME NONE' represents a lexeme that either hasn't been asigned yet
// or that it doesn't have a type.
const (
	NO_TYPE EntryType = iota
	INT
	FLOAT
	STRING
	FUNCTION
	BOOLEAN
)

func (l EntryType) Equals(a int) bool {
	return a == int(l)
}

func FromString(s string) EntryType {
	var str EntryType
	switch s {
	case "function":
		str = FUNCTION
	case "int":
		str = INT
	case "string":
		str = STRING
	case "float":
		str = FLOAT
	case "boolean":
		str = BOOLEAN
	}
	return str
}
func (l EntryType) String() string {
	var str string
	switch l {
	case FUNCTION:
		str = "function"
	case INT:
		str = "int"
	case STRING:
		str = "string"
	case FLOAT:
		str = "float"
	case BOOLEAN:
		str = "boolean"
	default:
		str = "indefinido"
	}
	return str
}

type Entry struct {
	id         int                   //id of Entry
	lexeme     string                //lexeme
	entry_type EntryType             //lexeme kind
	Attributes map[string]*Attribute //attribute list
	pos        int
}

func (e *Entry) GetPos() int {
	return e.pos
}

// Creates new symbol entry from the original lexeme.
// Returns the new Entry object with attributes init.
//
// Sets the type to LEXEME NONE by default
func NewEntry(lex string) *Entry {
	return &Entry{
		Attributes: map[string]*Attribute{},
		lexeme:     lex,
		entry_type: NO_TYPE,
	}

}
func (e *Entry) GetAttribute(a string) *Attribute {
	return e.Attributes[a]
}
func (e *Entry) GetType() EntryType {
	return e.entry_type
}

// sets the type of the lexeme
// IF an invalid lexem type is provided, an error is returned
func (e *Entry) setType(t EntryType, offset int) error {
	switch t {
	case FUNCTION:
	case INT:
	case STRING:
	case FLOAT:
	case BOOLEAN:
	case NO_TYPE:
	default:
		{
			if DEBUG {
				fmt.Printf("DEBUG: Invalid Lexeme Type: [%v]\n\r", t)
			}
			return fmt.Errorf("Error: Invalid Lexeme type: [%v]", t)
		}
	}
	e.entry_type = t
	e.SetAttributeValue("despl", offset)
	return nil
}

// Writes the SymbolEntry from the ST to the specified Writer with
// PDL specified format
func (e *Entry) Write() string {
	a := ""
	a += fmt.Sprintf("* LEXEMA: '%v'\r\n", e.lexeme)
	a += fmt.Sprintln("  Atributos:")
	a += "    + Tipo: "
	if e.entry_type == NO_TYPE {
		a += "'-'"
	} else {
		a += fmt.Sprintf("'%v'", e.entry_type.String())
	}
	a += fmt.Sprintln()
	for _, at := range e.Attributes {
		a += at.Write()
	}
	return a
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
			fmt.Print(k.Write())
		}
	}
}
