package errors

import (
	"fmt"
	"io"
)

var DEBUG bool
var manager ErrorManager

type ErrorManager struct {
	currentLine int
	errors      []Error
}

func NewErrorManager() {
	if DEBUG {
		fmt.Printf("DEBUG: Initializing Error Manager\n")
	}
	manager = ErrorManager{currentLine: 1, errors: []Error{}}
}
func Line()int{
	return manager.currentLine
}
// Create new Error.
func NewError(kind ErrorKind, code ErrorCode, val any) {
	manager.errors = append(manager.errors, newError(kind, code, manager.currentLine, val))
}
func SintacticalError(code ErrorCode, val any) {
	manager.errors = append(manager.errors, newError(SINTACTICAL, code, manager.currentLine, val))
}

func SemanticalError(code ErrorCode, val any) {
	manager.errors = append(manager.errors, newError(SEMANTICAL, code, manager.currentLine, val))
}
// When Lexer detecs new line, it calls this function to notifies the ErrorManager, so when the next error occurs,
// it can have the error's line
func NewLine() {
	manager.currentLine++
}

func Write(writer io.Writer) {
	if DEBUG && len(manager.errors) > 0 {
		fmt.Printf("DEBUG: Writing errors to output\n")
	}
	for _, e := range manager.errors {
		fmt.Fprintf(writer, "%v\n", e.string())
	}
}
