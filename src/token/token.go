package token

import (
	"bufio"
	"fmt"
)

type TokenKind int

var DEBUG bool

const (
	ABRIR_CORCH TokenKind = iota
	CERRAR_CORCH
	PUNTOYCOMA

	INT_LITERAL
	STRING_LITERAL
	REAL
	TRUE
	FALSE

	ABRIR_PAR
	CERRAR_PAR
	ARITM
	RELAC
	LOGICO

	ASIG
	ASIG_MULT
	ID
	LET

	INT
	FLOAT
	BOOLEAN
	STRING
	NULL

	WRITE
	READ

	DO
	WHILE
	FUNCTION
	RETURN
	VOID
	IF
	COMA
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
	case REAL:
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
	case ASIG_MULT:
		return "ASIG_MULT"
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
	Attr   string
}

func NewToken(kind TokenKind, lexeme string, attr string) Token {
	return Token{kind, lexeme, attr}
}

// Writes the token to the output with 'w' Writer.
func (tk *Token) Write(w *bufio.Writer) {
	if DEBUG {
		fmt.Fprintf(w, "<%v, %v>\n", tk.Kind.String(), tk.Attr)
	} else {
		fmt.Fprintf(w, "<%v, %v>\n", tk.Kind, tk.Attr)
	}
}
