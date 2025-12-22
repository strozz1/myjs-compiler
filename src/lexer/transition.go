package lexer

import (
	"compiler-pdl/src/errors"
	"compiler-pdl/src/token"
	"fmt"
	"math"
	"slices"
)

const MAX_STRING = 64

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
	errorState   State
}

func (t TransitionTable) isFinal() bool {
	return slices.Contains(t.finals, t.currentState)
}
func (t *TransitionTable) toError() {
	t.currentState = t.errorState
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

func (t *TransitionTable) Find(char rune) (*TransEntry, bool) {
	transition, ok := t.table[t.currentState]
	if !ok {
		return nil, false
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
		return nil, false
	}
	t.currentState = entry.Next
	return entry, true
}

// Generates de transitions of the DFA for the lexer.
func GenerateTransitions(sc *Lexer) TransitionTable {
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
		S11
		S12
		S13
		S14
		//Finals
		F0
		F1
		F2
		F3
		F4
		F5
		F6
		F7
		F8
		F9
		F10
		F11
		F12
		F13
		F14
		F15
		F16
		F17
		F18

		E1
		E2
	)
	finals := []State{
		F0, F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11,
		F12, F13, F14, F15, F16, F17, F18, E2}

	t := TransitionTable{
		table:        map[State]map[rune]*TransEntry{},
		finals:       finals,
		start:        S0,
		errorState:   E1,
		currentState: S0,
	}

	//delimiters
	t.addTransition(S0, S0, 0, matchDel, func() (token.Token, bool) {
		if sc.currentChar == '\n' {
			sc.newLine()
		}
		sc.nextChar()
		return token.Token{}, false
	})

	t.addTransition(S0, S1, '=', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	//equal
	t.addTransition(S1, F0, '=', nil, func() (token.Token, bool) {
		tk := token.NewToken(token.RELAC, "==", token.REL_EQ)
		sc.nextChar()
		return tk, true
	})

	//assign
	t.addTransition(S1, F1, 0, matchNotEq, func() (token.Token, bool) {
		tk := token.NewToken(token.ASIG, "=", token.ASIG_SIMPLE)
		return tk, true
	})

	t.addTransition(S0, S3, '*', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	//*=
	t.addTransition(S3, F2, '=', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.ASIG, "*=", token.ASIG_MULT)
		return tk, true
	})

	t.addTransition(S0, S4, '!', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	// !=
	t.addTransition(S4, F3, '=', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.RELAC, "!=", token.REL_NOTEQ)
		return tk, true
	})
	t.addTransition(S4, F4, 0, matchNotEq, func() (token.Token, bool) {
		tk := token.NewToken(token.LOGICO, "!", token.LOG_NEG)
		return tk, true
	})

	t.addTransition(S0, S5, -1, matchDigit, func() (token.Token, bool) {
		sc.appendChar()
		sc.intVal = 0
		d := int64(sc.currentChar - '0')
		sc.intVal += d
		sc.tokenState = NUMBER
		sc.nextChar()
		return token.Token{}, false
	})

	t.addTransition(S5, S5, -1, matchDigit, func() (token.Token, bool) {
		sc.appendChar()
		sc.intVal *= 10
		d := int64(sc.currentChar - '0')
		sc.intVal += d
		sc.nextChar()
		return token.Token{}, false
	})

	// END INT LITERAL
	t.addTransition(S5, F5, 0, matchNotDotOrDigit, func() (token.Token, bool) {
		value, ok := safeInt16(sc.intVal)
		if !ok {
			errors.NewError(errors.LEXICAL, errors.C_INT_TOO_BIG, sc.intVal)
			sc.reset()
			return token.Token{}, false
		}
		tk := token.NewToken(token.INT_LITERAL, fmt.Sprintf("%d", value), fmt.Sprintf("%d", value))
		return tk, true
	})

	t.addTransition(S5, S14, '.', nil, func() (token.Token, bool) {
		sc.appendChar()
		sc.decimalPos = 0
		sc.tokenState = FLOAT
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S14, S9, 0, matchDigit, func() (token.Token, bool) {
		sc.appendChar()
		sc.decimalPos += 1
		sc.intVal *= 10
		d := int64(sc.currentChar - '0')
		sc.intVal += d
		sc.nextChar()
		return token.Token{}, false
	})

	t.addTransition(S9, S9, 0, matchDigit, func() (token.Token, bool) {
		sc.decimalPos += 1
		sc.intVal *= 10
		d := int64(sc.currentChar - '0')
		sc.intVal += d
		sc.nextChar()
		return token.Token{}, false
	})
	//float
	t.addTransition(S9, F8, 0, matchNotDigit, func() (token.Token, bool) {
		var val float64
		if sc.intVal == 0 {
			val = 0
		} else {
			val = float64(sc.intVal) * math.Pow(10.0, float64(-sc.decimalPos))
		}
		value, ok := safeFloat32(val)
		if !ok {
			errors.NewError(errors.LEXICAL, errors.C_FLOAT_TOO_BIG, val)
			sc.reset()
			return token.Token{}, false
		}
		tk := token.NewToken(token.REAL_LITERAL, fmt.Sprintf("%f", value), fmt.Sprintf("%f", value))
		return tk, true
	})

	//BEGIN STRING_LITERAL
	t.addTransition(S0, S7, '\'', nil, func() (token.Token, bool) {
		sc.tokenState = STRING
		sc.nextChar()
		return token.Token{}, false
	})
	//CONT STRING LITERAL
	t.addTransition(S7, S7, 0, matchNotQuoteOrDel, func() (token.Token, bool) {
		sc.appendChar()
		sc.nextChar()
		return token.Token{}, false
	})
	//END STRING LITERAL
	t.addTransition(S7, F7, '\'', nil, func() (token.Token, bool) {
		if len(sc.lexeme) > MAX_STRING {
			errors.NewError(errors.LEXICAL, errors.C_STRING_TOO_LONG, sc.lexeme)
			sc.reset()
			return token.Token{}, false
		}
		sc.nextChar()
		tk := token.NewToken(token.STRING_LITERAL, sc.lexeme, sc.lexeme)
		return tk, true
	})

	// Start of ID
	t.addTransition(S0, S6, 0, matchIdFirst, func() (token.Token, bool) {
		sc.appendChar()
		sc.nextChar()
		return token.Token{}, false
	})

	// Cont of ID
	t.addTransition(S6, S6, 0, matchId, func() (token.Token, bool) {
		sc.appendChar()
		sc.nextChar()
		return token.Token{}, false
	})

	// END of ID
	t.addTransition(S6, F6, 0, matchEndId, func() (token.Token, bool) {
		var tk token.Token
		if sc.isReserved(sc.lexeme) {
			tk = token.NewToken(token.From(sc.lexeme), sc.lexeme, "")
		} else {
			if sc.declZone {
				val, ok := sc.STManager.AddEntry(sc.lexeme)
				tk = token.NewToken(token.ID, sc.lexeme, val)
				if ok {
					if DEBUG {
						fmt.Printf("DEBUG: Inserted new ID: %s with pos %d in ST\n", tk.Lexeme, tk.Attr.(int))
					}
				} else {
					errors.NewError(errors.SEMANTICAL, errors.SS_IDENTIFIER_DEFINED, sc.lexeme)
				}
			} else {
				entry, ok := sc.STManager.SearchEntry(sc.lexeme)
				if !ok {
					//global int
					val, _ := sc.STManager.AddGlobalEntry(sc.lexeme)
					e, _ := sc.STManager.GetEntry(val)
					sc.STManager.SetEntryType(e, "int")
					tk = token.NewToken(token.ID, sc.lexeme, val)
				} else {
					tk = token.NewToken(token.ID, sc.lexeme, entry.GetPos())
				}
			}
		}
		return tk, true
	})

	//comment
	t.addTransition(S0, S10, '/', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S10, S11, '*', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S11, S11, 0, matchNotStar, func() (token.Token, bool) {
		if sc.currentChar == '\n' {
			sc.newLine()
		}
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S11, S12, '*', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S12, S11, 0, matchNotDash, func() (token.Token, bool) {
		if sc.currentChar == '\n' {
			sc.newLine()
		}
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S12, S12, '*', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S12, S0, '/', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})

	//curl
	t.addTransition(S0, F10, '{', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.ABRIR_CORCH, "{", "")
		return tk, true
	})
	t.addTransition(S0, F11, '}', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.CERRAR_CORCH, "}", "")
		return tk, true
	})
	t.addTransition(S0, F12, '(', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.ABRIR_PAR, "(", "")
		return tk, true
	})
	t.addTransition(S0, F13, ')', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.CERRAR_PAR, ")", "")
		return tk, true
	})
	t.addTransition(S0, F14, ';', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.PUNTOYCOMA, ";", "")
		return tk, true
	})
	t.addTransition(S0, F15, '+', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.ARITM, "+", token.ARIT_PLUS)
		return tk, true
	})
	t.addTransition(S0, F16, '-', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.ARITM, "-", token.ARIT_MINUS)
		return tk, true
	})
	t.addTransition(S0, F17, ',', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.COMA, ",", "")
		return tk, true
	})

	t.addTransition(S0, S13, '&', nil, func() (token.Token, bool) {
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(S13, F18, '&', nil, func() (token.Token, bool) {
		sc.nextChar()
		tk := token.NewToken(token.LOGICO, "&&", token.LOG_AND)
		return tk, true
	})

	//error state
	t.addTransition(E1, E1, 0, matchNoDel, func() (token.Token, bool) {
		sc.appendChar()
		sc.nextChar()
		return token.Token{}, false
	})
	t.addTransition(E1, E2, 0, matchDelOrSemi, func() (token.Token, bool) {
		errors.NewError(errors.LEXICAL, sc.tokenState.toError(), sc.lexeme)
		return token.Token{}, false
	})

	//t.debugPrint()
	return t
}

func safeInt16(n int64) (int16, bool) {
	if n < math.MinInt16 || n > math.MaxInt16 {
		return 0, false
	}
	return int16(n), true
}

func safeFloat32(n float64) (float32, bool) {
	if n==0{return 0, true}
	if n < math.SmallestNonzeroFloat32 || n > math.MaxFloat32 {
		return 0, false
	}
	return float32(n), true
}

var matchNotDash = func(c rune) bool {
	return c != '/'
}
var matchNotStar = func(c rune) bool {
	return c != '*'
}
var matchNotQuoteOrDel = func(c rune) bool {
	return !(c == '\'' || c == '\n' || c == '\r')
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
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

var matchDelOrSemi = func(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == ';'
}
var matchNoDel = func(c rune) bool {
	return !(c == ' ' || c == '\t' || c == '\n' || c == '\r')
}
var matchNotEq = func(c rune) bool {
	return c != '='
}
var matchNotDotOrDigit = func(c rune) bool {
	return c != '.' && (c < '0' || c > '9')
}

var matchNotDigit = func(c rune) bool {
	return !(c >= '0' && c <= '9')
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
