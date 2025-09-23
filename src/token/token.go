package token

import (
	"bufio"
	"fmt"
)

type TokenKind int

var DEBUG bool

func From(token string) TokenKind {
	var t TokenKind
	switch token {
	case "do":
		t = DO
	case "while":
		t = WHILE
	case "var":
		t = VAR
	case "function":
		t = FUNCTION
	case "return":
		t = RETURN
	case "int":
		t = INT
	case "if":
		t = IF
	case "break":
		t = BREAK
	case "boolean":
		t = BOOLEAN
	case "NULL":
		t =NULL
	}
	return t
}

const (
	ID TokenKind = iota
	EOL
	ASIGNACION
	ENTERO_LIT
	CADENA_LIT
	ABRIR_PAR
	CERRAR_PAR
	ABRIR_CORCH
	CERRAR_CORCH
	COMA
	PUNTOYCOMA

	ARITM
	RELAC
	LOGICO

	//reserved
	DO
	WHILE
	VAR
	FUNCTION
	RETURN
	IF
	BREAK
	NULL
	BOOLEAN
	INT
)

func (t TokenKind) toString() string {
	var str string
	switch t {
	case EOL:
		str = "EOL"
	case CADENA_LIT:
		str = "CADENA_LIT"
	case ENTERO_LIT:
		str = "ENTERO_LIT"
	case ABRIR_CORCH:
		str = "ABRIR_CORCH"
	case CERRAR_CORCH:
		str = "CERRAR_CORCH"
	case ABRIR_PAR:
		str = "ABRIR_PAR"
	case CERRAR_PAR:
		str = "CERRAR_PAR"
	case ID:
		str = "ID"
	case DO:
		str = "DO"
	case WHILE:
		str = "WHILE"
	case ASIGNACION:
		str = "ASIGN"
	case VAR:
		str = "VAR"
	case FUNCTION:
		str = "FUNCTION"
	case RETURN:
		str = "RETURN"
	case COMA:
		str = "COMA"
	case PUNTOYCOMA:
		str = "PUNTOYCOMA"
	case ARITM:
		str = "ARITM"
	case LOGICO:
		str = "LOGICO"
	case RELAC:
		str = "RELAC"
	case IF:
		str = "IF"
	case BREAK:
		str = "BREAK"
	case BOOLEAN:
		str = "BOOLEAN"
	case NULL:
		str= "NULL"
	case INT:
		str = "INT"

	}
	return str
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
		fmt.Fprintf(w, "<%v, %v>\n", tk.Kind.toString(), tk.Attr)
	} else {
		fmt.Fprintf(w, "<%v, %v>\n", tk.Kind, tk.Attr)
	}
}
