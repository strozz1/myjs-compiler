package st

import (
	"fmt"
	"io"
)

var DEBUG bool

// ID counter for the tables
var stIdCounter = 0

type STManager struct {
	Global  *SymbolTable
	Local   *SymbolTable
	Current *SymbolTable

	ReservedWords []string
	//Defines the set of available attributes. This attributes are only the
	//template, so when you want to use it for an Entry you need to take it and modify it with the
	//corresponding values
	Attributes map[string]Attribute
	output     string
}

func (m *STManager) SearchEntry(lexeme string) (*Entry, bool) {
	a, ok := m.Current.table[lexeme]
	if !ok && (m.Current.name != m.Global.name) {
		a, ok = m.Global.table[lexeme]
	}
	return a, ok
}

func (m *STManager) EntryExists(e int) bool {
	_, b := m.GetEntry(e)
	return b
}

func (m *STManager) AddEntry(lexeme string) (int, bool) {
	return m.Current.AddEntry(lexeme)
}

// Creates a new SymbolTable Manager.
// Initializes the global SymbolTable.
func NewSTManager() *STManager {
	if DEBUG {
		fmt.Printf("DEBUG: Initializing STManager\n\r")
	}
	global := createST("Global Table")
	return &STManager{
		Global:     global,
		Attributes: map[string]Attribute{},
		Current:    global,
	}
}

// Creates an attribute with the specified name, type & description.
// 'AttributeType' is an int & 'AttributeDesc' is a string
func (m *STManager) CreateAttribute(name string, d string, t AttributeType) {
	if m.containsAttribute(name) {
		if DEBUG {
			fmt.Printf("DEBUG: Can't create attribute, it already exists [%v]\n\r", name)
			return
		}
	}
	m.Attributes[name] = NewAttribute(name, t, d)
	if DEBUG {
		fmt.Printf("DEBUG: Created new ST attribute '%v' of type: '%v' and description: '%v'\n\r", name, t, d)
	}
}

func (m *STManager) prepareFuncEntry(e *Entry) {

	a, _ := m.Attributes["tipoRetorno"]
	e.AddAtribute("tipoRetorno", a)
	a, _ = m.Attributes["numParam"]
	e.AddAtribute("numParam", a)
	a, _ = m.Attributes["etiqFuncion"]
	e.AddAtribute("etiqFuncion", a)
	e.SetAttributeValue("etiqFuncion", fmt.Sprintf("etiq_%v", e.lexeme))
}

func (m *STManager) SetEntryType(e *Entry, tt string) {
	t := FromString(tt)
	a, ok := m.Attributes["despl"]
	if !ok {
		if DEBUG {
			fmt.Printf("WARNING: 'despl' attribute not found\n")
		}
	} else {
		switch t {
		case FUNCTION:
			m.prepareFuncEntry(e)
		default:
			e.AddAtribute("despl", a)
		}
	}

	e.setType(t, m.Current.offset)
	m.shift(t)
}

func (m *STManager) shift(t EntryType) {
	var s int = 0
	switch t {
	case INT:
		s = 1
	case FLOAT:
		s = 2
	case STRING:
		s = 64
	case BOOLEAN:
		s = 1
	}
	m.Current.offset += s
}

func (m *STManager) GetEntry(pos int) (*Entry, bool) {
	e, ok := m.Current.GetEntry(pos)
	if !ok && (m.Current.name != m.Global.name) {
		e, ok = m.Global.GetEntry(pos)
	}
	return e, ok
}

func (m *STManager) GetGlobalEntry(pos int) (*Entry, bool) {
	return m.Global.GetEntry(pos)
}

func (m *STManager) SetEntryAttribute(e *Entry, name string, val any) {
	if !m.containsAttribute(name) {
		if DEBUG {
			fmt.Printf("DEBUG: Can't add attribute '%v'. It does not exist\n\r", name)
			return
		}
	}
	if !e.containsAttr(name) {
		v, _ := m.Attributes[name]
		e.AddAtribute(name, v)
		if DEBUG {
			fmt.Printf("DEBUG: Added attribute '%v' to entry '%v'\n\r", name, e.lexeme)
		}
	}

	e.SetAttributeValue(name, val)
	if DEBUG {
		fmt.Printf("DEBUG: Setted attribute value: '%v=%v' to entry\n\r", name, val)
	}
}

// @DEPRECATED
// Create a New local table with a specified name. If there was an existing local table
// it will be removed
func (m *STManager) CreateLocalTable(name string) {
	m.Local = createST(name)
	if DEBUG {
		fmt.Printf("DEBUG: Local table created: [%v]\n\r", name)
	}
}

// Creates a new scope with a 'name'.
// The new scope is the inner of the current scope and current will be parent
// of the new scope.
// Replaces current with new scope
func (m *STManager) NewScope(name string) {
	st := createST(name)
	st.parent = m.Current
	m.Current = st
	if DEBUG {
		fmt.Printf("DEBUG: New Scope '%v' created\n", name)
	}
}

// Destroy current Scope. This functions sets the current scope to the parent scope of 'current'.
// If 'm.Current' is Global Table operation is canceled.
func (m *STManager) DestroyScope() {
	if DEBUG {
		fmt.Printf("DEBUG: Scope '%v' destroyed\n", m.Current.name)
	}
	m.output = m.Current.Write() + m.output
	m.Current = m.Current.parent
}

// Returns if the attribute name already exists in the attribute list
func (m *STManager) containsAttribute(name string) bool {
	_, b := m.Attributes[name]
	return b
}

// Writes ST to the file specified
// @st: symbol table to write to the file
func (m *STManager) Write(writer io.Writer) {
	fmt.Fprintf(writer, "%s", m.output)
}
