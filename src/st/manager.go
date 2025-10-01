package st

import (
	"fmt"
	"io"
)

var DEBUG bool

// ID counter for the tables
var stIdCounter = 0

type STManager struct {
	Global        *SymbolTable
	Local         *SymbolTable
	Current       *SymbolTable

	ReservedWords []string
	//Defines the set of available attributes. This attributes are only the
	//template, so when you want to use it for an Entry you need to take it and modify it with the
	//corresponding values
	Attributes map[string]Attribute
}

// Creates a new SymbolTable Manager.
// Initializes the global SymbolTable.
func NewSTManager() *STManager {
	if DEBUG {
		fmt.Printf("DEBUG: Initializing STManager\n\r")
	}
	global:=createST("Global Table")
	return &STManager{
		Global: global,
		Attributes: map[string]Attribute{},
		Current: global,
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
		fmt.Printf("DEBUG: Created new attribute '%v' of type: '%v' & desc: '%v'\n\r", name, t, d)
	}
}

func (m *STManager) AddGlobalEntry(name string) *Entry {
	return m.Global.AddEntry(name)
}

func (m *STManager) AddLocalEntry(name string) *Entry {
	if m.Current == nil {
		if DEBUG {
			fmt.Println("DEBUG: Can't add entry to a 'nil' table.")
			return nil
		}
	}
	return m.Current.AddEntry(name)
}

func (m *STManager) GetGlobalEntry(name string) (*Entry, bool) {
	v, ok := m.Global.GetEntry(name)
	if !ok {
		if DEBUG {
			fmt.Printf("DEBUG: Global entry not found '%v'\n\r", name)
			return nil, ok
		}
	}
	return v, ok
}

func (m *STManager) RemoveGlobalEntry(name string) {
	m.Global.RemoveEntry(name)
}

func (m *STManager) RemoveLocalEntry(name string) {
	m.Current.RemoveEntry(name)
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
	st:=createST(name)
	m.Current.inner=st
	st.parent=m.Current
	m.Current=st
	if DEBUG{
		fmt.Printf("DEBUG: New Scope '%v' created\n",name)
	}
}

// Destroy current Scope. This functions sets the current scope to the parent scope of 'current'.
// If 'm.Current' is Global Table operation is canceled.
func (m *STManager) DestroyScope(){
	if m.Current==m.Global{
		if DEBUG{
			fmt.Printf("DEBUG: trying to destroy Global Table, operation canceled\n")
		}
	}else{
		if DEBUG{
			fmt.Printf("DEBUG: Scope '%v' destroyed\n",m.Current.name)
		}
		m.Current=m.Current.parent
	}
}

// Returns if the attribute name already exists in the attribute list
func (m *STManager) containsAttribute(name string) bool {
	_, b := m.Attributes[name]
	return b
}

// Writes ST to the file specified
// @st: symbol table to write to the file
func (m *STManager) Write(writer io.Writer, st *SymbolTable) {
	st.Write(writer)
}
