package diagnostic

import "fmt"

type ErrorKind int
const(
	ID_TOO_LONG ErrorKind = iota
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
	return fmt.Sprintf("ERROR en linea %v: %v",e.line,e.kind.string(e.info))
}

func (k ErrorKind) string(info string)string{
	var str string
	switch k{
	case ID_TOO_LONG:
		str=fmt.Sprintf("el identificador '%v' supera el limite de 128 caracteres",info)
	default:
		str="internal error"

	}
	return str
}
