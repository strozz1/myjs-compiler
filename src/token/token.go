package token

import (
	"bufio"
	"fmt"
)

var DEBUG bool

type TokenKind int

const (
	ID TokenKind = iota
	STRING_LITERAL
	INT_LITERAL
)

type Token struct {
	Kind   TokenKind
	Lexeme string
}

func NewToken(kind TokenKind, lexeme string) Token {
	return Token{kind, lexeme}
}

type TokenManager struct {
	tokens []Token
	writer *bufio.Writer
	tokenPos int
}

func (s *TokenManager) NextToken() (Token,bool){
	if len(s.tokens) == 0 || s.tokenPos >= len(s.tokens) {
		return Token{}, false
	}
	s.tokenPos++
	return s.tokens[s.tokenPos-1], true
}

func NewTokenManager(writer *bufio.Writer) *TokenManager {
	if DEBUG{
		fmt.Println("DEBUG: Initializing Token Manager")
	}
	return &TokenManager{
		tokens: []Token{},
		writer: writer,
	}
}

func (m *TokenManager) AddToken(tk Token) {

	m.tokens = append(m.tokens, tk)
	if DEBUG{
		fmt.Printf("DEBUG: Token added: <%v, %v>\n",tk.Kind,tk.Lexeme)
	}
}

func (m *TokenManager) Write(){
	for _,tk:=range m.tokens{
		fmt.Fprintf(m.writer,"<%v, %v>\n",tk.Kind,"-")
	}
	m.writer.Flush()
}
