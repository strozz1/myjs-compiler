package lexer

import (
	"bufio"
	"compiler-pdl/src/errors"
	"compiler-pdl/src/st"
	"compiler-pdl/src/token"
	"fmt"
	"io"
	"os"
	"slices"
)

var DEBUG bool

//TODO ver si error intchar
type TokenState int
const(
	NONE TokenState =iota
	NUMBER
	FLOAT
	ID
)
func (st TokenState) toError()errors.ErrorCode{
	var error errors.ErrorCode
	switch st{
	case NUMBER:
		error=errors.C_MALFORMED_NUMBER
	case FLOAT:
		error=errors.C_MALFORMED_FLOAT
	case ID:
		error=errors.C_MALFORMED_ID
	}
	return error
}

type Lexer struct {
	//currentChar char red
	currentChar rune
	//buffer red
	lexeme string

	intVal     int64
	decimalPos int
	tokenState TokenState

	//buffer input reader
	reader *bufio.Reader

	//transition table
	transitions TransitionTable

	//tokens
	tokens []token.Token

	//SymbolTable manager
	STManager *st.STManager

	EOF bool
}

// Creates a new lexer.
// Returns a lexer or error if failed to do so.
// The Lexer will be initialized and the first char will be red and
// saved in current
func NewLexer(r *bufio.Reader) (*Lexer, error) {
	char, _, err := r.ReadRune()
	if err != nil {
		return &Lexer{}, err
	}
	if DEBUG {
		fmt.Println("DEBUG: Initializing Lexer")
	}
	sc := Lexer{
		currentChar: char,
		reader:      r,
		tokens:      []token.Token{},
		STManager:   st.NewSTManager(),
	}
	sc.transitions = GenerateTransitions(&sc)
	return &sc, nil
}

// Check if 'token' is a reserved keyword
func (s *Lexer) isReserved(token string) bool {
	return slices.Contains(s.STManager.ReservedWords, token)
}

// Append current char to current token lexeme
func (s *Lexer) appendChar() {
	s.lexeme = string(append([]rune(s.lexeme), s.currentChar))
}

func (s *Lexer) newLine() {
	errors.NewLine()
}

// reads the next char from the input reader.
// Sets the value to the rune pointer
func (s *Lexer) nextChar() {
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

// The algorithm of the actual lexer. Saves the tokens in 's.tokens' and can be
// retreived one by one with 's.GetToken()'.
func (s *Lexer) Lexical() (token.Token, bool) {
	var ok bool = false
	var token token.Token
	for !ok && !s.EOF {
		transition, code, errVal := s.transitions.Find(s.currentChar)
		if transition == nil {
			if code ==-1 {
				if s.tokenState!=NONE{
					s.transitions.toError()
					continue
				}
				code=s.tokenState.toError()
				errVal=fmt.Sprintf("%s",s.lexeme)
			}else{
				code=errors.C_INVALID_CHAR
			}
			errors.NewError(errors.K_LEXICAL, code, errVal)
			s.reset()
			s.nextChar()
			return token, false //TODO
		}
		token, ok = transition.Action()
		if s.transitions.isFinal() {
			if ok {
				s.tokens = append(s.tokens, token)
			}
			s.reset()
			break
		}
	}
	return token, true
}

func (s *Lexer) reset() {
	s.lexeme = ""
	s.intVal = 0
	s.decimalPos = 0
	s.tokenState=0
	s.transitions.currentState = s.transitions.start
}

// WriteTokens all the tokens with the 'Writer' parameter. The output format follows the
// convention described in PDL subject.
func (s *Lexer) WriteTokens(w *bufio.Writer) {
	if DEBUG {
		fmt.Printf("DEBUG: Writting %v tokens to file\n", len(s.tokens))
	}
	for _, t := range s.tokens {
		t.Write(w)
	}
	w.Flush()
}

// Write lexical errors with the specified Writer
func (s *Lexer) WriteErrors(w io.Writer) {
	errors.Write(w)
}
