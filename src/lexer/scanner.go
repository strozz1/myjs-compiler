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
	intVal int32

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

	EOF bool
}

// Creates a new scanner.
// Returns a scanner or error if failed to do so.
// The Scanner will be initialized and the first char will be red and
// saved in current
func NewScanner(r *bufio.Reader) (*Scanner, error) {
	char, _, err := r.ReadRune()
	fmt.Printf("Red %c\n", char)
	if err != nil {
		return &Scanner{}, err
	}
	if DEBUG {
		fmt.Println("DEBUG: Initializing Lexer")
	}
	sc := Scanner{
		currentChar: char,
		reader:      r,
		tokens:      []token.Token{},
		STManager:   st.NewSTManager(),
		errManager:  diagnostic.NewErrorManager(),
	}
	sc.transitions = GenerateTransitions(&sc)
	return &sc, nil
}

// Check if 'token' is a reserved keyword
func (s *Scanner) isReserved(token string) bool {
	return slices.Contains(s.STManager.ReservedWords, token)
}

// Append current char to current token lexeme
func (s *Scanner) appendChar() {
	s.lexeme = string(append([]rune(s.lexeme), s.currentChar))
}

func (s *Scanner) newLine() {
	s.errManager.NewLine()
}

// reads the next char from the input reader.
// Sets the value to the rune pointer
func (s *Scanner) nextChar() {
	char, _, err := s.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			s.EOF = true
			if DEBUG {
				fmt.Println("DEBUG: EOF found. Finished reading input file")
			}
			return
		}
		if DEBUG {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		}
	}
	s.currentChar = char
}

// The algorithm of the actual scanner. Saves the tokens in 's.tokens' and can be
// retreived one by one with 's.GetToken()'.
func (s *Scanner) Lexical() (token.Token,bool) {
	var ok bool = false
	var token token.Token
	for !ok && !s.EOF {
		transition, code, errVal := s.transitions.Find(s.currentChar)
		if transition == nil {
			s.reset()
			s.errManager.NewError(diagnostic.K_LEXICAL, code, errVal)
			return token,false //TODO
		}
		token, ok = transition.Action()
		if s.transitions.isFinal() {
			s.tokens = append(s.tokens, token)
			s.reset()
			break
		}
	}
	return token,true
}

func (s *Scanner) reset() {
	s.lexeme = ""
	s.intVal = 0
	s.transitions.currentState = s.transitions.start
}

// Returns GetToken if there is a new one. This does not remove the token from the actual list
// Used to know what GetToken is next in the reading queue.
func (s *Scanner) GetToken() (token.Token, bool) {
	if len(s.tokens) == 0 || s.lastTokenRed >= len(s.tokens) {
		return token.Token{}, false
	}
	s.lastTokenRed++
	return s.tokens[s.lastTokenRed-1], true
}

// WriteTokens all the tokens with the 'Writer' parameter. The output format follows the
// convention described in PDL subject.
func (s *Scanner) WriteTokens(w *bufio.Writer) {
	if DEBUG {
		fmt.Printf("DEBUG: Writting %v tokens to file\n", len(s.tokens))
	}
	for _, t := range s.tokens {
		t.Write(w)
	}
	w.Flush()
}

// Write lexical errors with the specified Writer
func (s *Scanner) WriteErrors(w io.Writer) {
	s.errManager.Write(w)
}
