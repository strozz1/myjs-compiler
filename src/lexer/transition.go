package lexer

import (
	"compiler-pdl/src/token"
	"errors"
	"fmt"
)

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
	table   map[State]map[rune]*TransEntry
	initial State
	finals  []State
	current State
}

func (t *TransitionTable) addTransition(char rune, currentState State, nextState State, action Action, match Match) {
	trans, ok := t.table[currentState]
	if !ok {
		t.table[currentState] = map[rune]*TransEntry{}
		trans = t.table[currentState]
	}

	for r, i := range trans {
		//if a match set r
		if i.match != nil && i.match(char) {
			char = r
		}
	}
	_, ok = trans[char]
	if ok {
		if DEBUG {
			fmt.Printf("WARNING: Overwritting transition on ['S%v', %v]\n", currentState, char)
		}
	}
	trans[char] = &TransEntry{next: nextState, action: action, match: match}
	if DEBUG {
		fmt.Printf("DEBUG: New transition [S%v,char('%d')]=S%v]\n", currentState, char, nextState)
	}
}

func (t *TransitionTable) Find(char rune) (*TransEntry, error) {
	tr, ok := t.table[t.current]
	if !ok {
		return &TransEntry{}, errors.New("State not found")
	}
	for r, i := range tr {
		//if a match set r
		if i.match != nil && i.match(char) {

			char = r
		}
	}
	entry, ok := tr[char]
	if !ok {
		return entry, errors.New("Invalid char state")
	}
	t.current = entry.next
	return entry, nil
}

func GenerateTransitions(sc *Scanner) TransitionTable {
	t := TransitionTable{
		table:  map[State]map[rune]*TransEntry{},
		finals: []State{},
	}

	t.addTransition(' ', 0, 0, func(a rune, b *rune) {
		sc.Next(b)
	}, nil)
	t.addTransition('\n', 0, 0, func(a rune, b *rune) {
		sc.NewLine()
		sc.Next(b)
	}, nil)
	t.addTransition('\r', 0, 0, func(a rune, b *rune) {
		sc.Next(b)
	}, nil)
	t.addTransition(-2, 0, 1, func(a rune, b *rune) {
		sc.token = string(append([]rune(sc.token), a))
		sc.Next(b)
	}, matchID)

	t.addTransition(-3, 1, 1, func(a rune, b *rune) {
		sc.token = string(append([]rune(sc.token), a))
		sc.Next(b)
	}, matchID)
	t.addTransition(-4, 1, 0, func(a rune, b *rune) {
		var tk token.Token
		if !sc.IsReserved(sc.token) {
			tk = token.NewToken(token.ID, sc.token, "02")
			//TODO: check ST
		}else {
			//TODO que devolver en TIPO
			tk = token.NewToken(token.ID, sc.token, "02")

		}
		sc.AddToken(tk)
		t.current = t.initial
	}, matchDel)

	t.initial = 0
	t.current = 0
	t.debugPrint()
	return t
}

var matchID = func(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9')
}

// match delimiters
var matchDel = func(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == 0
}

func (m *TransitionTable) debugPrint() {
	if DEBUG {
		fmt.Printf("DEBUG: Printing transition table\n")
		for i, k := range m.table {
			fmt.Printf("Transitions for S%v\n", i)
			for j, w := range k {
				fmt.Printf("[%d]->%v\n", j, w.next)
			}
			fmt.Println("------------")
		}
	}
}
