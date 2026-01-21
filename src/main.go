package main

import (
	"bufio"
	"compiler-pdl/src/errors"
	"compiler-pdl/src/lexer"
	"compiler-pdl/src/parser"
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

	tkFile, e := os.Create("tokens.txt")
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error creating token file: %v\n", e)
	}
	defer tkFile.Close()
	stFile, e := os.Create("st.txt")
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error creating simbol table file: %v\n", e)
	}
	defer stFile.Close()
	parseFile, e := os.Create("parse.txt")
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error creating parse file: %v\n", e)
	}
	defer parseFile.Close()

	lexer, e := lexer.NewLexer(bufio.NewReader(file))
	initST(lexer.STManager)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error initializing lexer: %v\n", e)
		return
	}

	//init Error
	errors.NewErrorManager()

	parse := parser.NewParser(lexer)
	parse.Parse()
	lexer.STManager.Write(stFile)
	lexer.WriteTokens(bufio.NewWriter(tkFile))

	parse.Write(bufio.NewWriter(parseFile))
	lexer.WriteErrors(os.Stderr)
}



func initST(stManager *st.STManager) {
	stManager.ReservedWords = []string{
		"true", "false", "int", "float", "boolean", "string", "write", "read",
		"do", "while", "if", "function", "let", "return", "void",
	}
	stManager.CreateAttribute("despl", "despl", st.T_INTEGER)
	stManager.CreateAttribute("numParam", "numParam", st.T_INTEGER)
	stManager.CreateAttribute("tipoParam1", "tipoParam1", st.T_STRING)
	stManager.CreateAttribute("tipoParam2", "tipoParam2", st.T_STRING)
	stManager.CreateAttribute("tipoParam3", "tipoParam3", st.T_STRING)
	stManager.CreateAttribute("tipoParam4", "tipoParam4", st.T_STRING)
	stManager.CreateAttribute("tipoParam5", "tipoParam5", st.T_STRING)
	stManager.CreateAttribute("tipoParam6", "tipoParam6", st.T_STRING)
	stManager.CreateAttribute("tipoParam7", "tipoParam7", st.T_STRING)
	stManager.CreateAttribute("tipoParam8", "tipoParam8", st.T_STRING)
	stManager.CreateAttribute("tipoParam9", "tipoParam9", st.T_STRING)
	stManager.CreateAttribute("tipoRetorno", "tipoRetorno", st.T_STRING)
	stManager.CreateAttribute("etiqFuncion", "etiqFuncion", st.T_STRING)
	stManager.CreateAttribute("param", "param", st.T_INTEGER)
	stManager.CreateAttribute("dimension", "dimension", st.T_INTEGER)
}

func debug() {
	DEBUG = true
	errors.DEBUG = true
	token.DEBUG = true
	lexer.DEBUG = true
	parser.DEBUG = true
	st.DEBUG = true
}
