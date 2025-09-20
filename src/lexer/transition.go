package lexer

import (
	"compiler-pdl/src/diagnostic"
	"compiler-pdl/src/token"
	"errors"
	"fmt"
	"math"
)
const MAX_STRING=128
type State int
type Action func(rune, *rune)
type Match func(r rune) bool

// Defines an Entry for the transition table.
// consists on the Next State and a semantic Action
type TransEntry struct {
	next   State
	action Action
	match  Match
}

// Executes action for the transition
func (t *TransEntry) Action(a rune, b *rune) {
	t.action(a, b)
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
	trans[char] = &TransEntry{next: nextState, action: action, match: match}
	if DEBUG {
		fmt.Printf("DEBUG: New transition [S%v,char('%d')]=S%v]\n", currentState, char, nextState)
	}
}

func (t *TransitionTable) Find(char rune) (*TransEntry, error) {
	tr, ok := t.table[t.currentState]
	if !ok {
		return &TransEntry{}, errors.New("State not found")
	}
	for r, i := range tr {
		//if a match set r
		if i.match != nil && i.match(char) {
			char = r
			break
		}
	}
	entry, ok := tr[char]
	if !ok {
		return entry, fmt.Errorf("Caracter no valido '%c'",char)
	}
	t.currentState = entry.next
	return entry, nil
}

// Generates de transitions of the DFA for the lexer.
func GenerateTransitions(sc *Scanner) TransitionTable {
	t := TransitionTable{
		table:        map[State]map[rune]*TransEntry{},
		finals:       []State{},
		start:        0,
		currentState: 0,
	}
	const (
		S0 State = iota
		S1
		S2
		S3
		S4
		S5
		S6
		S7
	)

	//delimiters
	t.addTransition(S0, S0, 0, matchDel, func(a rune, b *rune) {
		sc.nextChar(b)
	})

	//new line
	t.addTransition(S0, S0, '\n', nil, func(a rune, b *rune) {
		sc.newLine()
		sc.nextChar(b)
	})
	//Carriage
	t.addTransition(S0, S0, '\r', nil, func(a rune, b *rune) {
		sc.nextChar(b)
	})

	// Start of ID
	t.addTransition(S0, S1, 0, matchIdFirst, func(a rune, b *rune) {
		sc.token = string(append([]rune(sc.token), a))
		sc.nextChar(b)
	})

	// Cont of ID
	t.addTransition(S1, S1, 0, matchId, func(a rune, b *rune) {
		sc.token = string(append([]rune(sc.token), a))
		sc.nextChar(b)
	})

	// END of ID
	t.addTransition(S1, S0, 0, matchDel, func(a rune, b *rune) {
		if !sc.IsReserved(sc.token) {
			sc.newToken(token.ID, "02")
			//TODO: check ST
		} else {
			//TODO que devolver en TIPO
			sc.newToken(token.ID, "01")
		}
	})

	//Start INT LITERAL
	t.addTransition(S0, S2, 0, matchDigit, func(a rune, next *rune) {
		sc.appendChar(a)
		sc.nextChar(next)
	})

	// Cont INT LITERAL
	t.addTransition(S2, S2, 0, matchDigit, func(a rune, b *rune) {
		sc.appendChar(a)
		sc.nextChar(b)
	})

	// END INT LITERAL
	t.addTransition(S2, S0, 0, matchEndInt, func(a rune, b *rune) {
		var raw int64 = 0
		var value int32
		for _, c := range sc.token {
			d := int64(c - '0')
			raw *= 10
			raw += d
			r, ok := safeInt32(raw)
			if !ok {
				sc.errManager.NewError(diagnostic.LEXICAL, fmt.Sprintf("el literal entero '%s' supera el maximo permitido.",sc.token))
				sc.token=""
				return
			}
			value=r
		}

		sc.newToken(token.INT_LITERAL, fmt.Sprintf("%d",value))
	})

	//BEGIN STRING_LITERAL
	t.addTransition(S0,S3,'\'',nil,func(a rune, b *rune){
		sc.appendChar(a)
		sc.nextChar(b)
	})
	//CONT STRING LITERAL
	t.addTransition(S3,S3,0,matchNotQuote,func(a rune,b *rune){
		sc.appendChar(a)
		sc.nextChar(b)
	})
	//END STRING LITERAL
	t.addTransition(S3,S0,'\'',nil,func(a rune,b *rune){
		sc.appendChar(a)
		if len(sc.token)-2 >MAX_STRING{
			sc.errManager.NewError(diagnostic.LEXICAL, fmt.Sprintf("La cadena literal %v supera el limite maximo de caracteres",sc.token))
		}else{
			sc.newToken(token.STRING_LITERAL,sc.token)
		}
		sc.nextChar(b)
	})

	//comment
	t.addTransition(S0,S4,'/',nil,func(a rune,b *rune){
		sc.nextChar(b)
	})
	t.addTransition(S4,S5,'*',nil,func(a rune, b *rune){
		sc.nextChar(b)
	})
	t.addTransition(S5,S6,0,matchNotStar,func(a rune, b *rune){
		if a=='\n'{
			sc.newLine()
		}
		sc.nextChar(b)
	})
	t.addTransition(S6,S6,0,matchNotStar,func(a rune, b *rune){
		if a=='\n'{
			sc.newLine()
		}
		sc.nextChar(b)
	})
	t.addTransition(S6,S7,'*',nil,func(a rune, b *rune){
		sc.nextChar(b)
	})
	t.addTransition(S7,S6,0,matchNotInv,func(a rune, b *rune){
		if a=='\n'{
			sc.newLine()
		}
		sc.nextChar(b)
	})
	t.addTransition(S7,S0,'/',nil,func(a rune, b *rune){
		sc.nextChar(b)
	})
	t.debugPrint()
	return t
}

func safeInt32(n int64) (int32, bool) {
	if n < math.MinInt32 || n > math.MaxInt32 {
		return 0, false
	}
	return int32(n), true
}
var matchNotInv= func(c rune)bool{
	return c!='/'
}
var matchNotStar = func(c rune)bool{
	return c!='*'
}
var matchNotQuote= func (c rune)bool{
	return c!='\''
}
var matchEndInt = func(c rune) bool {
	return !(c >= '0' && c <= '9')
}
var matchId = func(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9')
}
var matchIdFirst = func(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
var matchDigit = func(c rune) bool {
	return c >= '0' && c <= '9'
}

// match delimiters
var matchDel = func(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == 0
}

func (m *TransitionTable) debugPrint() {
	if DEBUG {
		fmt.Printf("DEBUG: Printing transition table\n")
		for i, k := range m.table {
			for j, w := range k {
				fmt.Printf("S%d(%d)->S%v\n", i, j, w.next)
			}
			fmt.Println("------------")
		}
	}
}
