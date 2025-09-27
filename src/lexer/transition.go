package lexer

import (
	"compiler-pdl/src/diagnostic"
	"compiler-pdl/src/token"
	"errors"
	"fmt"
	"math"
	"slices"
)

const MAX_STRING = 128

// Defines a Match for the transition
type Match func(r rune) bool

// Defines a Semantical Action for the lexer to perform
type Action func() (token.Token, bool)
type State int

// Defines an Entry for the transition table.
// consists on the Next State and a Semantic Action and optionaly a match. This function,
// if present will be used to check if the transition is valid. We use this to not generate N
// transition for all the possible letter of the alfabet corresponding the transition
type TransEntry struct {
	Next   State
	Action Action
	Match  Match
}

// Reprensents the transition table from the DFA
type TransitionTable struct {
	table        map[State]map[rune]*TransEntry
	start        State
	finals       []State
	currentState State
}

var transId rune = -1

// Creates a new transition from currentState to nextState.
// It can transition either using 1 char or via a Match. If match is NOT nil, the char will be ignored.
// Internally, if a match is present, the table stores the transition with a negative char(that autodecrements)
func (t *TransitionTable) addTransition(currentState State, nextState State, char rune, match Match, action Action) {
	trans, ok := t.table[currentState]
	if !ok {
		t.table[currentState] = map[rune]*TransEntry{}
		trans = t.table[currentState]
	}

	if match != nil {
		char = transId
		transId--
	}
	_, ok = trans[char]

	if ok {
		if DEBUG {
			fmt.Printf("WARNING: Overwritting transition on ['S%v', %v]=S%d\n", currentState, char, nextState)
		}
	}
	trans[char] = &TransEntry{Next: nextState, Action: action, Match: match}
}

func (t *TransitionTable) Find(char rune) (*TransEntry, error) {
	//check if we are in a final state
	if slices.Contains(t.finals, t.currentState) {
		t.currentState = t.start
	}
	transition, ok := t.table[t.currentState]
	if !ok {
		return &TransEntry{}, errors.New("State not found")
	}
	for r, i := range transition {
		//if a match set r
		if i.Match != nil && i.Match(char) {
			char = r
			break
		}
	}
	entry, ok := transition[char]
	if !ok {
		return entry, fmt.Errorf("Caracter no valido '%c'", char)
	}
	t.currentState = entry.Next
	return entry, nil
}

// Generates de transitions of the DFA for the lexer.
func GenerateTransitions(sc *Scanner) TransitionTable {
	const (
		S0 State = iota
		S1
		S2
		S3
		S4
		S5
		S6
		S7
		S8
		S9
		S10

		F1
		F2
		F3
	)
	finals := []State{F1, F2, F3}
	t := TransitionTable{
		table:        map[State]map[rune]*TransEntry{},
		finals:       finals,
		start:        S0,
		currentState: S0,
	}

	//delimiters
	t.addTransition(S0, S0, 0, matchDel, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})

	//new line
	t.addTransition(S0, S0, '\n', nil, func() (token.Token, bool) {
		sc.newLine()
		sc.nextChar()
		return token.Token{}, false
	})
	//Carriage
	t.addTransition(S0, S0, '\r', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})

	// Start of ID
	t.addTransition(S0, S1, 0, matchIdFirst, func() (token.Token, bool) {
		sc.lexeme = string(append([]rune(sc.lexeme), sc.currentChar))
		sc.nextChar()
		return token.Token{}, false
	})

	// Cont of ID
	t.addTransition(S1, S1, 0, matchId, func() (token.Token, bool) {
		sc.lexeme = string(append([]rune(sc.lexeme), sc.currentChar))
		sc.nextChar()
		return token.Token{}, false
	})

	// END of ID
	t.addTransition(S1, S0, 0, matchEndId, func() (token.Token, bool) {
		var tk token.Token
		if sc.isReserved(sc.lexeme) {
			tk = token.NewToken(token.ID, sc.lexeme, "-")
		} else {
			tk = token.NewToken(token.ID, sc.lexeme, "-")
			//TODO: check ST
		}
		return tk, true
	})

	//Start INT LITERAL
	t.addTransition(S0, S2, 0, matchDigit, func() (token.Token, bool) {
		sc.intVal *= 10
		d := int32(sc.currentChar - '0')
		sc.intVal += d
		sc.nextChar()
		return token.Token{}, false
	})

	// Cont INT LITERAL
	t.addTransition(S2, S2, 0, matchDigit, func() (token.Token, bool) {
		sc.intVal *= 10
		d := int32(sc.currentChar - '0')
		sc.intVal += d
		sc.nextChar()
		return token.Token{}, false
	})

	// END INT LITERAL
	t.addTransition(S2, S0, 0, matchEndInt, func() (token.Token, bool) {
		value, ok := safeInt16(sc.intVal)
		if !ok {
			sc.errManager.NewError(diagnostic.LEXICAL, fmt.Sprintf("el literal entero '%d' supera el maximo permitido.", sc.intVal))
			sc.lexeme = ""
			return token.Token{}, false
		}
		tk := token.NewToken(token.INT_LITERAL, "", fmt.Sprintf("%d", value))
		return tk, true
	})

	//BEGIN STRING_LITERAL
	t.addTransition(S0, S3, '\'', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	//CONT STRING LITERAL
	t.addTransition(S3, S3, 0, matchNotQuote, func() (token.Token, bool) {
		sc.appendChar()
		sc.nextChar()
		return token.Token{}, false
	})
	//END STRING LITERAL
	t.addTransition(S3, S0, '\'', nil, func() (token.Token, bool) {
		if len(sc.lexeme) > MAX_STRING {
			sc.errManager.NewError(diagnostic.LEXICAL, fmt.Sprintf("La cadena literal %v supera el limite maximo de caracteres", sc.lexeme))
			return token.Token{}, false
		}
		tk:=token.NewToken(token.STRING_LITERAL,sc.lexeme,sc.lexeme)
		sc.nextChar()
		return tk,true
	})

	//comment
	t.addTransition(S0, S4, '/', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S4, S5, '*', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S5, S6, 0, matchNotStar, func() (token.Token, bool) {
		if sc.currentChar == '\n' {
			sc.newLine()
		}
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S6, S6, 0, matchNotStar, func() (token.Token, bool) {
		if sc.currentChar == '\n' {
			sc.newLine()
		}
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S6, S7, '*', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S7, S6, 0, matchNotInv, func() (token.Token, bool) {
		if sc.currentChar == '\n' {
			sc.newLine()
		}
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S7, S0, '/', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})

	//curl
	t.addTransition(S0, S0, '{', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.ABRIR_CORCH,"{","-")
		sc.nextChar()
		return tk, true
	})
	t.addTransition(S0, S0, '}', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.CERRAR_CORCH,"}","-")
		sc.nextChar()
		return tk, true
	})
	t.addTransition(S0, S0, '(', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.ABRIR_PAR,"(","-")
		sc.nextChar()
		return tk, true
	})
	t.addTransition(S0, S0, ')', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.CERRAR_PAR,")","-")
		sc.nextChar()
		return tk, true
	})
	t.addTransition(S0, S0, ',', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.COMA,",","-")
		sc.nextChar()
		return tk, true
	})
	t.addTransition(S0, S0, ';', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.PUNTOYCOMA,";","-")
		sc.nextChar()
		return tk, true
	})

	//asign
	t.addTransition(S0, S0, '=', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{},false
	})

	t.addTransition(S0, S0, '+', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.ARITM,"+","&")
		//TODO tipo
		sc.nextChar()
		return tk, true
	})
	t.addTransition(S0, S8, '&', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S8, S0, '&', nil, func() (token.Token, bool) {
		tk:=token.NewToken(token.LOGICO,"&&","&&")
		sc.nextChar()
		return tk, true
	})
	//t.debugPrint()
	return t
}

func safeInt16(n int32) (int16, bool) {
	if n < math.MinInt16 || n > math.MaxInt16 {
		return 0, false
	}
	return int16(n), true
}

var matchNotInv = func(c rune) bool {
	return c != '/'
}
var matchNotStar = func(c rune) bool {
	return c != '*'
}
var matchNotQuote = func(c rune) bool {
	return c != '\''
}
var matchEndInt = func(c rune) bool {
	return !(c >= '0' && c <= '9')
}
var matchId = func(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9')
}
var matchEndId = func(c rune) bool {
	return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9'))
}
var matchIdFirst = func(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
var matchDigit = func(c rune) bool {
	return c >= '0' && c <= '9'
}

// match delimiters
var matchDel = func(c rune) bool {
	return c == ' ' || c == '\t' || c == 0
}

func (m *TransitionTable) debugPrint() {
	if DEBUG {
		fmt.Printf("DEBUG: Printing transition table\n")
		for i, k := range m.table {
			for j, w := range k {
				fmt.Printf("S%d(%d)->S%v\n", i, j, w.Next)
			}
			fmt.Println("------------")
		}
	}
}
