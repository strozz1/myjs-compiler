package st

import (
	"os"
)

func main() {
	DEBUG=true
	stManager := NewSTManager()
	stManager.CreateAttribute("despl", DESC_DESPL, T_INTEGER)
	stManager.CreateAttribute("numero de parametros", DESC_NUM_PARAM, T_INTEGER)
	stManager.CreateAttribute("tipo de parametros", DESC_TIPO_PARAM, T_ARRAY)
	stManager.CreateAttribute("modo de parametros", DESC_MODO_PARAM, T_ARRAY)
	stManager.CreateAttribute("tipo de retorno", DESC_TIPO_RETORNO, T_STRING)
	stManager.CreateAttribute("etiqueta", DESC_ETIQ_FUNCION,T_STRING)
	stManager.CreateAttribute("parametro", DESC_PARAM,T_INTEGER)
	stManager.CreateAttribute("dimension", "dimension",T_INTEGER)
	stManager.CreateAttribute("elem", "elem",T_STRING)


	var1 := stManager.AddGlobalEntry("x")
	var1.SetType(LEXEME_INTEGER)
	stManager.SetEntryAttribute(var1, "despl", 0)
	stManager.SetEntryAttribute(var1, "parametro", 0)
	
	var2 := stManager.AddGlobalEntry("y")
	var2.SetType(LEXEME_INTEGER)
	stManager.SetEntryAttribute(var2, "despl", 4)
	stManager.SetEntryAttribute(var2, "parametro", 0)

	var3 := stManager.AddGlobalEntry("res")
	var3.SetType(LEXEME_INTEGER)
	stManager.SetEntryAttribute(var3,"despl",8)
	stManager.SetEntryAttribute(var3,"parametro",0)
	
	func1:=stManager.AddGlobalEntry("suma")
	func1.SetType(LEXEME_FUNCTION)
	stManager.SetEntryAttribute(func1,"numero de parametros",2)
	tipos:=[]string{"int","int"}
	stManager.SetEntryAttribute(func1,"tipo de parametros",tipos)
	modos:=[]string{"0","0"}
	stManager.SetEntryAttribute(func1,"modo de parametros",modos)
	stManager.SetEntryAttribute(func1,"tipo de retorno","int")
	stManager.SetEntryAttribute(func1,"etiqueta","sumatorio")

	stManager.CreateLocalTable("suma")

	local1:=stManager.AddLocalEntry("a")
	local1.SetType(LEXEME_INTEGER)
	stManager.SetEntryAttribute(local1,"despl",0)
	stManager.SetEntryAttribute(local1,"parametro",1)

	local2:=stManager.AddLocalEntry("b")
	local2.SetType(LEXEME_INTEGER)

	stManager.SetEntryAttribute(local2,"despl",4)
	stManager.SetEntryAttribute(local2,"parametro",1)

	local3:=stManager.AddLocalEntry("res")
	local3.SetType(LEXEME_INTEGER)

	stManager.SetEntryAttribute(local3,"despl",8)
	stManager.SetEntryAttribute(local3,"parametro",0)


	stManager.Write(os.Stdout,stManager.Local)
	stManager.Write(os.Stdout,stManager.Global)
}
