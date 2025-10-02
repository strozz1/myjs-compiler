package token

import (
	"bufio"
	"fmt"
)

type TokenKind int

var DEBUG bool

const (
	EOF TokenKind = iota
	BOOLEAN
	DO
	FLOAT
	FUNCTION
	IF
	INT
	LET
	READ
	RETURN
	STRING
	VOID
	WHILE
	WRITE
	REAL_LITERAL
	INT_LITERAL
	STRING_LITERAL
	ID
	ASIG
	COMA
	PUNTOYCOMA
	ABRIR_PAR
	CERRAR_PAR
	ABRIR_CORCH
	CERRAR_CORCH
	ARITM
	LOGICO
	RELAC

	FALSE
	TRUE

	NULL
)

type AsigType int

const (
	ASIG_MULT AsigType = iota
	ASIG_SIMPLE 
)

type AritType int

const (
	ARIT_PLUS AritType = iota
	ARIT_MINUS
)

type LogicType int

const (
	LOG_AND LogicType = iota
	LOG_NEG 
)
type RelType int
const(
	REL_NOTEQ LogicType=iota
	REL_EQ
)

func (t TokenKind) String() string {
	switch t {
	case ABRIR_CORCH:
		return "ABRIR_CORCH"
	case CERRAR_CORCH:
		return "CERRAR_CORCH"
	case PUNTOYCOMA:
		return "PUNTOYCOMA"
	case INT_LITERAL:
		return "INT_LITERAL"
	case STRING_LITERAL:
		return "STRING_LITERAL"
	case REAL_LITERAL:
		return "REAL"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case ABRIR_PAR:
		return "ABRIR_PAR"
	case CERRAR_PAR:
		return "CERRAR_PAR"
	case ARITM:
		return "ARITM"
	case RELAC:
		return "RELAC"
	case LOGICO:
		return "LOGICO"
	case ASIG:
		return "ASIG"
	case ID:
		return "ID"
	case LET:
		return "LET"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case BOOLEAN:
		return "BOOLEAN"
	case STRING:
		return "STRING"
	case NULL:
		return "NULL"
	case WRITE:
		return "WRITE"
	case READ:
		return "READ"
	case DO:
		return "DO"
	case WHILE:
		return "WHILE"
	case FUNCTION:
		return "FUNCTION"
	case RETURN:
		return "RETURN"
	case VOID:
		return "VOID"
	case IF:
		return "IF"
	case COMA:
		return "COMA"
	default:
		return "UNKNOWN"
	}
}

func From(token string) TokenKind {
	var t TokenKind
	switch token {
	case "true":
		t = TRUE
	case "false":
		t = FALSE
	case "let":
		t = LET
	case "int":
		t = INT
	case "float":
		t = FLOAT
	case "string":
		t = STRING
	case "boolean":
		t = BOOLEAN
	case "null":
		t = NULL
	case "write":
		t = WRITE
	case "read":
		t = READ
	case "do":
		t = DO
	case "while":
		t = WHILE
	case "function":
		t = FUNCTION
	case "return":
		t = RETURN
	case "if":
		t = IF
	case "void":
		t = VOID
	}
	return t
}

type Token struct {
	Kind   TokenKind
	Lexeme string
	Attr   any
}

func NewToken(kind TokenKind, lexeme string, attr any) Token {
	return Token{kind, lexeme, attr}
}

// Writes the token to the output with 'w' Writer.
func (tk *Token) Write(w *bufio.Writer) {
	if DEBUG {
		fmt.Fprintf(w, "<%v, %v>\n", tk.Kind.String(), tk.Attr)
	} else {
		if tk.Kind==STRING_LITERAL{
			fmt.Fprintf(w, "<%d, \"%v\">\n", tk.Kind, tk.Attr)
		}else{
			fmt.Fprintf(w, "<%d, %v>\n", tk.Kind, tk.Attr)
		}
	}
}
