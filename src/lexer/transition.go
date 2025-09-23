package lexer

import (
	"compiler-pdl/src/diagnostic"
	"compiler-pdl/src/token"
	"errors"
	"fmt"
	"math"
)
const MAX_STRING=128

//Defines a Match for the transition
type Match func(r rune) bool
//Defines a Semantical Action for the lexer to perform
type Action func()
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
	tr, ok := t.table[t.currentState]
	if !ok {
		return &TransEntry{}, errors.New("State not found")
	}
	for r, i := range tr {
		//if a match set r
		if i.Match != nil && i.Match(char) {
			char = r
			break
		}
	}
	entry, ok := tr[char]
	if !ok {
		return entry, fmt.Errorf("Caracter no valido '%c'",char)
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
	)
	t := TransitionTable{
		table:        map[State]map[rune]*TransEntry{},
		finals:       []State{},
		start:        S0,
		currentState: S0,
	}

	//delimiters
	t.addTransition(S0, S0, 0, matchDel, func() {
		sc.nextChar()
	})

	//new line
	t.addTransition(S0, S0, '\n', nil, func() {
		sc.newLine()
		sc.nextChar()
	})
	//Carriage
	t.addTransition(S0, S0, '\r', nil, func() {
		sc.nextChar()
	})

	// Start of ID
	t.addTransition(S0, S1, 0, matchIdFirst, func() {
		sc.lexeme = string(append([]rune(sc.lexeme), sc.currentChar))
		sc.nextChar()
	})

	// Cont of ID
	t.addTransition(S1, S1, 0, matchId, func() {
		sc.lexeme = string(append([]rune(sc.lexeme), sc.currentChar))
		sc.nextChar()
	})

	// END of ID
	t.addTransition(S1, S0, 0, matchEndId, func() {
		if sc.isReserved(sc.lexeme) {
			sc.newToken(token.From(sc.lexeme), "-")
		} else {
			sc.newToken(token.ID, "-")
			//TODO: check ST
		}
	})

	//Start INT LITERAL
	t.addTransition(S0, S2, 0, matchDigit, func() {
		sc.appendChar()
		sc.nextChar()
	})

	// Cont INT LITERAL
	t.addTransition(S2, S2, 0, matchDigit, func() {
		sc.appendChar()
		sc.nextChar()
	})

	// END INT LITERAL
	t.addTransition(S2, S0, 0, matchEndInt, func() {
		var raw int64 = 0
		var value int32
		for _, c := range sc.lexeme {
			d := int64(c - '0')
			raw *= 10
			raw += d
			r, ok := safeInt32(raw)
			if !ok {
				sc.errManager.NewError(diagnostic.LEXICAL, fmt.Sprintf("el literal entero '%s' supera el maximo permitido.",sc.lexeme))
				sc.lexeme=""
				return
			}
			value=r
		}

		sc.newToken(token.INT_LITERAL, fmt.Sprintf("%d",value))
	})

	//BEGIN STRING_LITERAL
	t.addTransition(S0,S3,'\'',nil,func(){
		sc.appendChar()
		sc.nextChar()
	})
	//CONT STRING LITERAL
	t.addTransition(S3,S3,0,matchNotQuote,func(){
		sc.appendChar()
		sc.nextChar()
	})
	//END STRING LITERAL
	t.addTransition(S3,S0,'\'',nil,func(){
		sc.appendChar()
		if len(sc.lexeme)-2 >MAX_STRING{
			sc.errManager.NewError(diagnostic.LEXICAL, fmt.Sprintf("La cadena literal %v supera el limite maximo de caracteres",sc.lexeme))
		}else{
			sc.newToken(token.STRING_LITERAL,sc.lexeme)
		}
		sc.nextChar()
	})

	//comment
	t.addTransition(S0,S4,'/',nil,func(){
		sc.nextChar()
	})
	t.addTransition(S4,S5,'*',nil,func(){
		sc.nextChar()
	})
	t.addTransition(S5,S6,0,matchNotStar,func(){
		if sc.currentChar=='\n'{
			sc.newLine()
		}
		sc.nextChar()
	})
	t.addTransition(S6,S6,0,matchNotStar,func(){
		if sc.currentChar=='\n'{
			sc.newLine()
		}
		sc.nextChar()
	})
	t.addTransition(S6,S7,'*',nil,func(){
		sc.nextChar()
	})
	t.addTransition(S7,S6,0,matchNotInv,func(){
		if sc.currentChar=='\n'{
			sc.newLine()
		}
		sc.nextChar()
	})
	t.addTransition(S7,S0,'/',nil,func(){
		sc.nextChar()
	})

	//curl
	t.addTransition(S0,S0,'{',nil,func(){
		sc.newToken(token.OPEN_CURLY,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S0,'}',nil,func(){
		sc.newToken(token.CLOSE_CURLY,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S0,'(',nil,func(){
		sc.newToken(token.OPEN_PAR,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S0,')',nil,func(){
		sc.newToken(token.CLOSE_PAR,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S0,',',nil,func(){
		sc.newToken(token.COMMA,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S0,';',nil,func(){
		sc.newToken(token.SEMICOLON,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S0,'=',nil,func(){
		sc.newToken(token.ASIGN,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S0,'+',nil,func(){
		sc.newToken(token.PLUS,"-")
		sc.nextChar()
	})
	t.addTransition(S0,S8,'&',nil,func(){
		sc.nextChar()
	})
	t.addTransition(S8,S0,'&',nil,func(){
		sc.newToken(token.LOGIC_AND,"-")
		sc.nextChar()
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
var matchEndId=func(c rune) bool{
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
