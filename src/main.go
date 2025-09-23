package main

import (
	"bufio"
	"compiler-pdl/src/diagnostic"
	"compiler-pdl/src/lexer"
	"compiler-pdl/src/st"
	"compiler-pdl/src/token"
	"fmt"
	"os"
)

var DEBUG bool

func main() {
	_, ok := os.LookupEnv("DEBUG")
	if ok {
		debug()
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "You must specify a file.\n")
		os.Exit(-1)
	}
	path := os.Args[1]
	file, e := os.Open(path)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", e)
		return
	}
	defer file.Close()

	tkFile, e := os.Create("output/tokens.txt")
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error creating token file: %v\n", e)
	}
	defer tkFile.Close()

	scanner, e := lexer.NewScanner(bufio.NewReader(file))
	initST(scanner.STManager)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error initializing lexer: %v\n", e)
		return
	}

	scanner.ScanTokens()
	scanner.WriteTokens(bufio.NewWriter(tkFile))
	scanner.WriteErrors(os.Stderr)
}

func initST(stManager *st.STManager) {
	stManager.ReservedWords = []string{
		"do", "while", "if", "function", "var", "return", "break",
	}
	stManager.CreateAttribute("despl", "despl", st.T_INTEGER)
	stManager.CreateAttribute("numero de parametros", "numParam", st.T_INTEGER)
	stManager.CreateAttribute("tipo de parametros", "tipoParam", st.T_ARRAY)
	stManager.CreateAttribute("modo de parametros", "modoParam", st.T_ARRAY)
	stManager.CreateAttribute("tipo de retorno", "tipoRetorno", st.T_STRING)
	stManager.CreateAttribute("etiqueta", "etiqFuncion", st.T_STRING)
	stManager.CreateAttribute("parametro", "param", st.T_INTEGER)
	stManager.CreateAttribute("dimension", "dimension", st.T_INTEGER)
	stManager.CreateAttribute("elem", "elem", st.T_STRING)

}

func debug() {
	DEBUG = true
	diagnostic.DEBUG = true
	lexer.DEBUG = true
	token.DEBUG = true
	st.DEBUG = true
}
