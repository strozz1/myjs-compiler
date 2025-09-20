package token

import (
	"bufio"
	"fmt"
)

var DEBUG bool

type TokenKind int

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

// Manage the existing tokens. This manager does NOT create new tokens or detects them.
// The purpose of this manager is to
type TokenManager struct {
	tokens   []Token
	writer   *bufio.Writer
	tokenPos int
	token    string
}

func NewTokenManager(writer *bufio.Writer) *TokenManager {
	if DEBUG {
		fmt.Println("DEBUG: Initializing Token Manager")
	}
	return &TokenManager{
		tokens: []Token{},
		writer: writer,
	}
}

// Returns Token if there is a new one. This does not remove the token from the actual list
// Used to know what Token is next in the reading queue.
func (s *TokenManager) PopToken() (Token, bool) {
	if len(s.tokens) == 0 || s.tokenPos >= len(s.tokens) {
		return Token{}, false
	}
	s.tokenPos++
	return s.tokens[s.tokenPos-1], true
}

// Push a new token
func (m *TokenManager) PushToken(tk Token) {
	m.tokens = append(m.tokens, tk)
	if DEBUG {
		fmt.Printf("DEBUG: Token added: <%v, %v>\n", tk.Kind.toString(), tk.Attr)
	}
}

func (m *TokenManager) Write() {
	for _, tk := range m.tokens {
		if DEBUG {
			fmt.Fprintf(m.writer, "<%v, %v>\n", tk.Kind.toString(), tk.Attr)
		} else {
			fmt.Fprintf(m.writer, "<%v, %v>\n", tk.Kind, tk.Attr)
		}
	}
	m.writer.Flush()
}
