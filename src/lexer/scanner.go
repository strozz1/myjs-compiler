package lexer

import (
	"bufio"
	"compiler-pdl/src/diagnostic"
	"compiler-pdl/src/st"
	"compiler-pdl/src/token"
	"fmt"
	"io"
	"os"
	"slices"
)

var DEBUG bool

type Scanner struct {
	//currentChar char red
	currentChar rune
	//buffer red
	lexeme string
	intVal int

	//buffer input reader
	reader *bufio.Reader

	//transition table
	transitions TransitionTable
	//tokens
	tokens []token.Token

	//last token red from outsiders(not last token created)
	lastTokenRed int

	//SymbolTable manager
	STManager *st.STManager
	//Error Manager
	errManager diagnostic.ErrorManager
}


// Creates new token with the current buffer string.
// resets the buffer 'lexeme' & 'intVal' for next tokens.
func (sc *Scanner) newToken(d token.TokenKind, attr string) {
	tk := token.NewToken(d, sc.lexeme, attr)
	sc.tokens = append(sc.tokens, tk)
	sc.lexeme = ""
	sc.intVal = 0
}

// Check if it's a reserved keyword
func (s *Scanner) isReserved(token string) bool {
	return slices.Contains(s.STManager.ReservedWords, token)
}

// Append char to current token lexeme
func (s *Scanner) appendChar(c rune) {
	s.lexeme = string(append([]rune(s.lexeme), c))
}

// WriteTokens all the tokens with the 'Writer' parameter. The output format follows the
// convention described in PDL subject.
func (s *Scanner) WriteTokens(w *bufio.Writer) {
	if DEBUG {
		fmt.Printf("DEBUG: Writting %v tokens to file\n",len(s.tokens))
	}
	for _, t := range s.tokens {
		t.Write(w)
	}
	w.Flush()
}

//Write lexical errors with the specified Writer
func (s *Scanner) WriteErrors(w io.Writer){
	s.errManager.Write(w)
}

func (s *Scanner) newLine() {
	s.errManager.NewLine()
}

// Creates a new scanner for the reader provided.
// debug sets the debug mode.
// Returns a scanner or error if failed to do so.
//
// The Scanner will be initialized and the first char will be red and
// saved in current
func NewScanner(r *bufio.Reader) (*Scanner, error) {
	char, _, err := r.ReadRune()
	fmt.Printf("Red %c\n",char)
	if err != nil {
		return &Scanner{}, err
	}
	if DEBUG {
		fmt.Println("DEBUG: Initializing Lexer")
	}
	sc := Scanner{
		currentChar:    char,
		reader:     r,
		tokens:     []token.Token{},
		STManager:         st.NewSTManager(),
		errManager: diagnostic.NewErrorManager(),
	}
	sc.transitions = GenerateTransitions(&sc)
	return &sc, nil
}

// reads the next char from the input reader.
// Sets the value to the rune pointer
func (s *Scanner) nextChar(next *rune) {
	char, _, err := s.reader.ReadRune()
	*next = char
	if err != nil {
		if err == io.EOF {
			if DEBUG {
				fmt.Println("DEBUG: EOF found. Finished reading input file")
			}
			return
		}
		if DEBUG {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		}
		return
	}
}

// The algorithm of the actual scanner
func (s *Scanner) ScanTokens() {
	for s.currentChar != 0 {
		e, err := s.transitions.Find(s.currentChar)
		if err != nil {
			s.errManager.NewError(diagnostic.LEXICAL, err.Error())
			return //TODO
		}
		e.action(s.currentChar, &s.currentChar)
	}
}

// Returns Token if there is a new one. This does not remove the token from the actual list
// Used to know what Token is next in the reading queue.
func (s *Scanner) Token() (token.Token, bool) {
	if len(s.tokens) == 0 || s.lastTokenRed >= len(s.tokens) {
		return token.Token{}, false
	}
	s.lastTokenRed++
	return s.tokens[s.lastTokenRed-1], true
}
