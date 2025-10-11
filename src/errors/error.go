package errors

import "fmt"

type ErrorKind int

const (
	K_SINTACTICAL ErrorKind = iota
	K_LEXICAL
	K_SEMANTICAL
)

type ErrorCode int

const (
	C_OK ErrorCode = iota
	C_INVALID_CHAR
	C_STRING_TOO_LONG 
	C_INT_TOO_BIG
	C_FLOAT_TOO_BIG
	C_MALFORMED_NUMBER
	C_MALFORMED_FLOAT
	C_MALFORMED_STRING
	C_MALFORMED_ID
)

type Error struct {
	kind ErrorKind
	code ErrorCode
	line int
	val  any
}

// Creates new error from kind and line number
func newError(kind ErrorKind, code ErrorCode, line int, val any) Error {
	return Error{kind, code, line, val}
}

func (e *Error) string() string {
	return fmt.Sprintf("ERROR %v(%d) en la LINEA %v: %v", e.kind.string(), e.code, e.line, e.code.string(e.val))
}

func (k ErrorKind) string() string {
	var str string
	switch k {
	case K_LEXICAL:
		str = "LEXICO"
	case K_SINTACTICAL:
		str = "SINTACTICO"
	case K_SEMANTICAL:
		str = "SEMANTICO"
	default:
		str = "INTERNO"
	}
	return str
}

func (c ErrorCode) string(val any) string {
	var str string
	switch c {
	case C_STRING_TOO_LONG:
		str = fmt.Sprintf("La cadena '%s' supera el limite maximo de caracteres", val)
	case C_INVALID_CHAR:
		str = fmt.Sprintf("Caracter invalido: '%c'", val)
	case C_INT_TOO_BIG:
		str = fmt.Sprintf("Valor entero '%d' supera el limite.", val)
	case C_FLOAT_TOO_BIG:
		str = fmt.Sprintf("Valor real '%f' supera el limite.", val)
	case C_MALFORMED_NUMBER:
		str= fmt.Sprintf("Literal numerico '%s' mal formado",val)
	case C_MALFORMED_FLOAT:
		str= fmt.Sprintf("Literal real '%s' mal formado",val)
	case C_MALFORMED_ID:
		str= fmt.Sprintf("El identificador '%s' no es valido",val)
	case C_MALFORMED_STRING:
		str= fmt.Sprintf("Cadena '%s' mal formada",val)
	case C_OK:
		str = "OK"
	default:
		str = "interal error"
	}
	return str
}
