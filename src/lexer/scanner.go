package lexer

import (
	"bufio"
	"compiler-pdl/src/diagnostic"
	"compiler-pdl/src/st"
	"compiler-pdl/src/token"
	"fmt"
)

var DEBUG bool

type Scanner struct {
	//current char red
	current rune
	//buffer reader
	reader *bufio.Reader
	//transition table
	transitions TransitionTable
	//token list
	tokens []token.Token
	// Token list position
	tokenPos int
	//SymbolTable manager
	st *st.STManager
	//Error Manager
	errManager *diagnostic.ErrorManager
}

// Creates a new scanner for the reader provided.
// debug sets the debug mode.
// Returns a scanner or error if failed to do so.
//
// The Scanner will be initialized and the first char will be red and
// saved in current
func NewScanner(r *bufio.Reader, st *st.STManager, diagnostic *diagnostic.ErrorManager) (Scanner, error) {

	char, _, err := r.ReadRune()
	if err != nil {
		return Scanner{}, err
	}
	if DEBUG {
		fmt.Println("DEBUG: Initializing Lexer")
	}
	return Scanner{
		current:    char,
		reader:     r,
		tokens:     []token.Token{},
		tokenPos:   0,
		st:         st,
		errManager: diagnostic,
		transitions: GenerateTransitions(),
	}, nil
}

// Returns true if next input exists.
func (s *Scanner) hasNext() bool {
	return s.current != 0
}

// The algorithm of the actual scanner
func (s *Scanner) Scan() {
	for s.hasNext() {

	}
}

// Returns the last token generated. If token found also returns true.
// Internally, it checks s.tokenPos(last token gotten) and if new tokens returns it.
func (s *Scanner) Token() (token.Token, bool) {
	if len(s.tokens) == 0 || s.tokenPos >= len(s.tokens) {
		return token.Token{}, false
	}
	s.tokenPos++
	return s.tokens[s.tokenPos-1], true
}
