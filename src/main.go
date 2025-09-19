package main

import (
	"bufio"
	"compiler-pdl/src/diagnostic"
	"compiler-pdl/src/lexer"
	"compiler-pdl/src/token"
	"fmt"
	"compiler-pdl/src/st"
	"os"
)

var DEBUG bool
func main(){
	_,er:=os.LookupEnv("DEBUG")
	if er{
		debug()
	}

	if(len(os.Args)!= 2){
		fmt.Fprintf(os.Stderr,"You must specify a file.\n")
		os.Exit(-1)
	}
	path:=os.Args[1]
	file,e:=os.Open(path)
	if e!=nil{
		fmt.Fprintf(os.Stderr,"Error reading file: %v\n",e)
		return
	}
	symTable:=st.CreateSTManager(os.Stdout);
	errManager:=diagnostic.NewErrorManager(os.Stderr)

	lexer,e:=lexer.NewScanner(bufio.NewReader(file),&symTable,&errManager)
	if e!=nil{
		fmt.Fprintf(os.Stderr,"Error initializing lexer: %v\n",e)
		return
	}
	lexer.Scan()
}

func debug() {
	DEBUG=true
	diagnostic.DEBUG=true
	lexer.DEBUG=true
	token.DEBUG=true
	st.DEBUG=true
}
