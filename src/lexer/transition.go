package lexer

import (
	"compiler-pdl/src/token"
	"errors"
	"fmt"
)

type State int
type Action func(rune,*rune)
// Defines an Entry for the transition table.
// consists on the Next State and a semantic Action
type TransEntry struct{
	next State
	action Action
}

// Executes action for the transition
func (t *TransEntry) Action(a rune, b *rune){
	t.action(a,b)
}

// Reprensents the transition table from the DFA
type TransitionTable struct{
	table map[rune]map[State]*TransEntry
	initial State
	finals []State
	current State
}

func (t *TransitionTable) addTransition(char rune,currentState State,nextState State, action Action){
	trans,ok:=t.table[char]
	if !ok{
		t.table[char]=map[State]*TransEntry{}
		trans=t.table[char]
	}
	_,ok=trans[currentState]
	if ok{
		if DEBUG{
			fmt.Printf("WARNING: Overwritting transition on ['%v',S%v]\n",char,currentState)
		}
	} 
	trans[currentState]=&TransEntry{next:nextState,action:action}
	if DEBUG{
		fmt.Printf("DEBUG: New transition ['%c',S%v]=S%v]\n",char,currentState,nextState)
	}
}

func (t *TransitionTable) Find(char rune) (*TransEntry,error){
	tr,ok:=t.table[char]
	if !ok{
		return &TransEntry{}, errors.New("Char not supported")
	}
	entry,ok:=tr[t.current]
	if !ok{
		return entry, errors.New("Invalid char state")
	}
	t.current=entry.next
	return entry,nil
}


func GenerateTransitions(sc *Scanner)TransitionTable{
	t:=TransitionTable{
		table: map[rune]map[State]*TransEntry{},
		finals: []State{},
	}
	
	t.addTransition(' ',0,0,func(a rune, b *rune){
		sc.Next(b)
	})
	t.addTransition('\n',0,0,func(a rune, b *rune){
		sc.NewLine()
		sc.Next(b)
	})
	t.addTransition('\r',0,0,func(a rune, b *rune){
		sc.Next(b)
	})
	t.addTransition('v',0,0,func(a rune, b *rune){
		token:=token.NewToken(token.STRING_LITERAL,"v")
		sc.AddToken(token)
		t.current=t.initial
		sc.Next(b)
	})
	t.initial=0
	t.current=0

	return t
}
