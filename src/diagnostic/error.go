package diagnostic

import "fmt"

type ErrorKind int
const(
	SINTACTICAL ErrorKind = iota
	LEXICAL
)

type Error struct{
	kind ErrorKind
	line int
	info string
}

//Creates new error from kind and line number
func NewError(kind ErrorKind,line int,info string)Error{
	return Error{kind,line,info}
}

func (e *Error) ToString()string{
	return fmt.Sprintf("ERROR %v en linea %v: %v",e.kind.string(),e.line,e.info)
}

func (k ErrorKind) string()string{
	var str string
	switch k{
	case LEXICAL:
		str="LEXICO"
	case SINTACTICAL:
		str="SINTACTICO"
	default:
		str="DESCONOCIDO"
	}
	return str
}
