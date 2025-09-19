package lexer

import "fmt"

type State int
type Action func(rune)
// Defines an Entry for the transition table.
// consists on the Next State and a semantic Action
type TransEntry struct{
	next State
	action Action
}

// Executes action for the transition
func (t *TransEntry) Action(r rune){
	t.action(r)
}

// Reprensents the transition table from the DFA
type TransitionTable struct{
	table map[rune]map[State]TransEntry
}

func (t *TransitionTable) addTransition(char rune,currentState State,nextState State, action Action){
	trans,ok:=t.table[char]
	if !ok{
		t.table[char]=map[State]TransEntry{}
		trans=t.table[char]
	}
	_,ok=trans[currentState]
	if ok{
		if DEBUG{
			fmt.Printf("WARNING: Overwritting transition on ['%v',S%v]\n",char,currentState)
		}
	} 
	trans[currentState]=TransEntry{next:nextState,action:action}
	if DEBUG{
		fmt.Printf("DEBUG: New transition ['%c',S%v]=S%v]\n",char,currentState,nextState)
	}
}

func GenerateTransitions()TransitionTable{
	
	t:=TransitionTable{
		table: map[rune]map[State]TransEntry{},
	}

	t.addTransition('v',0,1,func(rune){
		
	})


	return t
}
