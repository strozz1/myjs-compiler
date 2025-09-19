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
	T_ARRAY
	T_NONE
)

// PREDEFINED ATTRIBUTES
const (
	DESC_DESPL        = "Despl"       // relative offset
	DESC_NUM_PARAM    = "numParam"    //num of params
	DESC_TIPO_PARAM   = "TipoParam"   // type of params
	DESC_MODO_PARAM   = "ModoParam"   // param mode
	DESC_TIPO_RETORNO = "TipoRetorno" // return type
	DESC_ETIQ_FUNCION = "EtiqFuncion" // function label
	DESC_PARAM        = "Param"       // param
)

type Attribute struct {
	Name      string        // att name
	Type      AttributeType //type of attribute
	Desc      string        //Description of attribute
	StringVal string        //String value if Type is String
	IntVal    int           //Int value if Type is Int
	ArrayVal  []string      //Array value if Type is Array
	HasValue  bool          //Flag if value has been asigned
}

// Creates new attribute
func NewAttribute(name string, tp AttributeType, ad string) Attribute {
	return Attribute{
		Name:     name,
		Type:     tp,
		Desc:     ad,
		HasValue: false,
	}
}

// Writes the Attribute from an Entry to the specified Writer with
// PDL specified format
func (a *Attribute) Write(w io.Writer) {
	switch a.Type {
	case T_INTEGER:
		fmt.Fprintf(w, "    + %v: %v\n\r", a.Desc, a.IntVal)
	case T_STRING:
		if a.HasValue {
			fmt.Fprintf(w, "    + %v: '%v'\n\r", a.Desc, a.StringVal)
		} else {
			fmt.Fprintf(w, "    + %v: '-'\n\r", a.Desc)
		}
	case T_ARRAY:
		if a.HasValue {
			for i, v := range a.ArrayVal {
				fmt.Fprintf(w, "    + %v%v: '%v'\n\r", a.Desc, i, v)
			}
			return
		} else {
			fmt.Fprintf(w, "    + %v: '-'\n\r", a.Desc)
		}
	}
}
