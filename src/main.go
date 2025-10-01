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

	lexer, e := lexer.NewLexer(bufio.NewReader(file))
	initST(lexer.STManager)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error initializing lexer: %v\n", e)
		return
	}

	for !lexer.EOF{
		lexer.Lexical()
	}
	lexer.WriteTokens(bufio.NewWriter(tkFile))
	lexer.WriteErrors(os.Stderr)
}

// PREDEFINED ATTRIBUTES
const (
	DESC_DESPL        = "Despl"       // relative offset
	DESC_NUM_PARAM    = "numParam"    //num of params
	DESC_TIPO_PARAM   = "TipoParam"   // type of params
	DESC_MODO_PARAM   = "ModoParam"   // param mode
	DESC_TIPO_RETORNO = "TipoRetorno" // return type
	DESC_ETIQ_FUNCION = "EtiqFuncion" // function label
	DESC_PARAM        = "Param"       // param
)


func initST(stManager *st.STManager) {
	stManager.ReservedWords = []string{
		"true","false","int","float","boolean","string","write","read",
		"do", "while", "if", "function", "let", "return", "void",
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
