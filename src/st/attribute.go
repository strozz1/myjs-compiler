package st

import (
	"fmt"
	"io"
)

// Defines the types that an attribute can be
type AttributeType int

const (
	T_INTEGER AttributeType = iota
	T_STRING
	T_BOOLEAN
	T_ARRAY
	T_NONE
)

type Attribute struct {
	Name      string        // att name
	Type      AttributeType //type of attribute
	Desc      string        //Description of attribute
	stringVal string        //String value if Type is String
	intVal    int           //Int value if Type is Int
	arrayVal  []string      //Array value if Type is Array
	hasValue  bool          //Flag if value has been asigned
}


func (a *Attribute) Value() any {
	switch a.Type {
	case T_INTEGER:
		return a.intVal
	case T_STRING:
		return a.stringVal
	case T_ARRAY:
		return a.arrayVal
	}
	return nil
}

// Creates new attribute
func NewAttribute(name string, tp AttributeType, ad string) Attribute {
	return Attribute{
		Name:     name,
		Type:     tp,
		Desc:     ad,
		hasValue: false,
	}
}

// Writes the Attribute from an Entry to the specified Writer with
// PDL specified format
func (a *Attribute) Write(w io.Writer) {
	switch a.Type {
	case T_INTEGER:
		fmt.Fprintf(w, "    + %v: %v\n\r", a.Desc, a.intVal)
	case T_STRING:
		if a.hasValue {
			fmt.Fprintf(w, "    + %v: '%v'\n\r", a.Desc, a.stringVal)
		} else {
			fmt.Fprintf(w, "    + %v: '-'\n\r", a.Desc)
		}
	case T_ARRAY:
		if a.hasValue {
			for i, v := range a.arrayVal {
				fmt.Fprintf(w, "    + %v%v: '%v'\n\r", a.Desc, i, v)
			}
			return
		} else {
			fmt.Fprintf(w, "    + %v: '-'\n\r", a.Desc)
		}
	}
}
