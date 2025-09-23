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
	}
	return t
}

const (
	ID TokenKind = iota
	STRING_LITERAL
	INT_LITERAL
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAR
	CLOSE_PAR
	DO
	WHILE
	ASIGN
	VAR
	FUNCTION
	RETURN
	COMMA
	SEMICOLON
	PLUS
	LOGIC_AND
)

func (t TokenKind) toString() string {
	var str string
	switch t {
	case STRING_LITERAL:
		str = "STRING_LITERAL"
	case INT_LITERAL:
		str = "INT_LITERAL"
	case OPEN_CURLY:
		str = "OPEN_CURLY"
	case CLOSE_CURLY:
		str = "CLOSE_CURLY"
	case OPEN_PAR:
		str = "OPEN_PAR"
	case CLOSE_PAR:
		str = "CLOSE_PAR"
	case ID:
		str = "ID"
	case DO:
		str = "DO"
	case WHILE:
		str = "WHILE"
	case ASIGN:
		str = "ASIGN"
	case VAR:
		str = "VAR"
	case FUNCTION:
		str = "FUNCTION"
	case RETURN:
		str = "RETURN"
	case COMMA:
		str = "COMMA"
	case SEMICOLON:
		str = "SEMICOLON"
	case PLUS:
		str = "PLUS"
	case LOGIC_AND:
		str = "LOGIC_AND"
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
