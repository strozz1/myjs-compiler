package errors

import "fmt"

type ErrorKind int

const (
	SINTACTICAL ErrorKind = iota
	LEXICAL
	SEMANTICAL
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

	S_INVALID_EXP
	S_EXPECTED_EXP
	S_EXPECTED_WHILE_CORCH
	S_EXPECTED_SENT
	S_EXPECTED_CERRAR_CORCH
	S_MISSING_WHILE
	S_EXPECTED_ABRIR_PAR
	S_EXPECTED_CERRAR_PAR
	S_EXPECTED_FUNCTYPE
	S_EXPECTED_SEMICOLON
	S_EXPECTED_RET_EXP
	S_EXPECTED_EXP_LOG
	S_TYPE

	SS_IDENTIFIER_DEFINED
	SS_ID_NOT_FOUND
	SS_INVALID_EXP_TYPE
	SS_NUM_PARAMS_INV
	SS_EXPECTED_WHILE_COND
	SS_EXPECTED_BOOLEANS
	SS_NEGATION_EXPECTED_BOOL
	SS_INVALID_ARIT_TYPES
	SS_RELATIONAL_TYPES
SS_INVALID_RETURN
	SS_RETURN_OUTSIDE
	SS_INVALID_SIGN_TYPE
SS_IF_COND
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
	case LEXICAL:
		str = "LEXICO"
	case SINTACTICAL:
		str = "SINTACTICO"
	case SEMANTICAL:
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
		str = fmt.Sprintf("Literal numerico '%s' mal formado", val)
	case C_MALFORMED_FLOAT:
		str = fmt.Sprintf("Literal real '%s' mal formado", val)
	case C_MALFORMED_ID:
		str = fmt.Sprintf("El identificador '%s' no es valido", val)
	case C_MALFORMED_STRING:
		str = fmt.Sprintf("Cadena '%s' mal formada", val)
	case C_OK:
		str = "OK"
		//SINTACTICAL
	case S_INVALID_EXP:
		str = "La expresion no es valida"
	case S_TYPE:
		str = "Se espera un Tipo"
	case S_EXPECTED_EXP:
		str = "Se esperaba expresion"
	case S_EXPECTED_WHILE_CORCH:
		str = "Se espera { antes del cuerpo del while"
	case S_EXPECTED_SENT:
		str = "Se esperaba sentencia"
	case S_EXPECTED_CERRAR_CORCH:
		str = "Falta el cierre de bloque '}'"
	case S_MISSING_WHILE:
		str = "Se esperaba keyword 'while'"
	case S_EXPECTED_CERRAR_PAR:
		str = "Se esperaba cierre de parentesis ')'"
	case S_EXPECTED_ABRIR_PAR:
		str = "Se esperaba apertura de parentesis '('"
	case S_EXPECTED_SEMICOLON:
		str = "Falta ';' al final de la expresion"
	case S_EXPECTED_FUNCTYPE:
		str = "La funcion no tiene definida el tipo de retorno"
	case S_EXPECTED_RET_EXP:
		str = "Se esperaba expresion de retorno"
	case S_EXPECTED_EXP_LOG:
		str = "Se esperaba expresion logica"
	case SS_IDENTIFIER_DEFINED:
		str = fmt.Sprintf("El identificador '%s' ya existe", val)
	case SS_ID_NOT_FOUND:
		str = fmt.Sprintf("El identificador '%s' no existe", val)
	case SS_INVALID_EXP_TYPE:
		str = fmt.Sprintf("Se espera expresion de tipo '%s'", val)
	case SS_NUM_PARAMS_INV:
		str = fmt.Sprintf("Numero de parametros invalidos: %s", val)
	case SS_EXPECTED_WHILE_COND:
		str = fmt.Sprintf("Condicion do-while erronea. Se esperaba 'boolean', se tiene '%s'", val)
	case SS_EXPECTED_BOOLEANS:
		str = fmt.Sprintf("Se esperaba comparacion de dos expresiones logicas, se tiene '%s'", val)
	case SS_NEGATION_EXPECTED_BOOL:
		str = fmt.Sprintf("Se esperaba un boolean en la operacion logica de negacion, se obtuvo '%s'", val)
	case SS_INVALID_ARIT_TYPES:
		str = fmt.Sprintf("Tipos aritmeticos invalidos, se esperaban tipos numericos iguales, se obtuvo %s", val)
	case SS_INVALID_SIGN_TYPE:
		str = fmt.Sprintf("No se puede usar signo al tipo '%s'", val)
	case SS_RELATIONAL_TYPES:
		str = fmt.Sprintf("Los tipos de la expresion de relacion deben ser los mismos, se obtuvo %s", val)
	case SS_RETURN_OUTSIDE:
		str = "Expresion de retorno fuera de funcion"
	case SS_INVALID_RETURN:
		str = fmt.Sprintf("Tipo de retorno invalido, %s", val)
	case SS_IF_COND:
		str = fmt.Sprintf("La condicion del IF debe ser booleana, se obtuvo '%s'", val)
	default:
		str = "interal error"
	}
	return str
}
