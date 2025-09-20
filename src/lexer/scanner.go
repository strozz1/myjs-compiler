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
	//current char red
	current rune
	//buffer red
	token string
	//buffer reader
	reader *bufio.Reader
	//transition table
	transitions TransitionTable
	//token manager
	tkManager *token.TokenManager
	//SymbolTable manager
	st *st.STManager
	//Error Manager
	errManager *diagnostic.ErrorManager
}

// Creates new token with the current buffer string.
// It calls the token Manager to Push the new token and
// resets the buffer 'token' for next tokens.
func (sc *Scanner) newToken(d token.TokenKind, param2 string) {
	tk := token.NewToken(d, sc.token, param2)
	sc.tkManager.PushToken(tk)
	sc.token = ""
}

// Check if it's a reserved keyword
func (s *Scanner) IsReserved(token string) bool {
	return slices.Contains(s.st.ReservedWords, token)
}

//Append char to current token lexeme
func (s *Scanner) appendChar(c rune) {
	s.token = string(append([]rune(s.token), c))
}
func (s *Scanner) Write() {
	if DEBUG {
		fmt.Printf("DEBUG: Writting tokens to file\n")
	}
	s.tkManager.Write()
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
func NewScanner(r *bufio.Reader, st *st.STManager, diagnostic *diagnostic.ErrorManager, tkManager *token.TokenManager) (Scanner, error) {
	char, _, err := r.ReadRune()
	if err != nil {
		return Scanner{}, err
	}
	if DEBUG {
		fmt.Println("DEBUG: Initializing Lexer")
	}
	sc := Scanner{
		current:    char,
		reader:     r,
		tkManager:  tkManager,
		st:         st,
		errManager: diagnostic,
	}
	sc.transitions = GenerateTransitions(&sc)
	return sc, nil
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
func (s *Scanner) Scan() {
	for s.current != 0 {
		e, err := s.transitions.Find(s.current)
		if err != nil {
			s.errManager.NewError(diagnostic.LEXICAL, err.Error())
			return //TODO
		}
		e.action(s.current, &s.current)
	}
}

// Returns the last token generated. If token found also returns true.
// Internally, it checks s.tokenPos(last token gotten) and if new tokens returns it.
func (s *Scanner) Token() (token.Token, bool) {
	return s.tkManager.PopToken()
}
