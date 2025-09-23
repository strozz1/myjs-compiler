package diagnostic

import (
	"fmt"
	"io"
)

var DEBUG bool

type ErrorManager struct{
	currentLine int
	errors []Error
}
func NewErrorManager()ErrorManager{
	if DEBUG{
		fmt.Printf("DEBUG: Initializing Error Manager\n")
	}
	return ErrorManager{currentLine: 1,errors:[]Error{}}

}

func (m *ErrorManager) NewError(kind ErrorKind,info string){
	m.errors = append(m.errors, NewError(kind,m.currentLine,info))
}

// When Lexer detecs new line, it calls this function to notifies the ErrorManager, so when the next error occurs,
// it can have the error's line
func (m *ErrorManager) NewLine(){
	m.currentLine++
}

func (m *ErrorManager) Write(writer io.Writer){
	if DEBUG && len(m.errors)>0{
		fmt.Printf("DEBUG: Writing errors to output\n")
	}
	for _,e:=range m.errors{
		fmt.Fprintf(writer,"%v\n",e.ToString())
	}
}
