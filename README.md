# Practica Compilador MyJS
Este proyecto esta dedicado a aprender a como desarrollar un compilador para un lenguaje sencillo llamado `MyJS`.

La especificacion completa del lenguaje se puede encontrar en la web del departamento de lenguajes de la [ETSIINF](https://dlsiis.fi.upm.es/procesadores/Practica.html).


El proyecto esta divido en varios modulos:
- Analizador Lexico
- Analizador Sintactico
- Analizador Semantico
- ...

En cada modulo se explicaran todos los detalles y especificaciones necesarias para cada parte.

## Especificacion del lenguaje
El lenguaje de este compilador es una version muy reducida del lenguaje `JavaScript`, llamada **MyJS**. Esta version tiene lo minimo que
necesita un lenguaje para poder programar y ejecutarlo en una maquina.


## Analizador Lexico
## Analizador Sintantico
Para el analizador sintactico se disponen distintas tecnicas de parsear lo que nos devuelve el Lexico.

Se nos ha seleccionado el metodo de analisis descendente recursivo, o **Recursive Descent LL(1)**. Este metodo sintactico analiza el programa fuente de izquierda a derecha de arriba a abajo. Es decir, expandiendo las producciones de la gramatica (al contrario que el ascendente, que las reduce).

Una gramatica G debe tener unas propiedades para que forme un lenguaje $L(G)$ apto para $LL(1)$:
- **Factorizada**: debe estar factorizada
- **No recursiva por la izquierda**: esto haria que entrase en bucle infinito
- **No ambigua**: una gramatica ambigua generaria un conflicto a la hora de elegir las posibles producciones.


A continuacion se detalla la gramatica creada para dicho lenguaje teniendo en cuenta las restricciones y especificaciones del lenguaje.
### Gramatica LL(1)
Una produccion es de la forma $A->\alpha_1|alpha_2$ donde $A$ es la parte izquierda de la produccion o cabeza, y $\alpha_1$, $\alpha_2$ la parte derecha o consecuente.

La parte izquierda solo puede contener un No-Terminal, y la derecha cualquier combinacion de terminales y no terminales.

Se define el conjunto de No Terminales:
```c
ReturnExp Sent2 Sent ParamList2 ParamList FuncBody FuncParams2 FuncParams TipoFunc DecFunc Tipo FactorId Term2 Term3 Term AritExp2 AritExp ExpRel2 ExpRel Expr2 Expr WhileBody TipoDecl Decl P
```

Se define el conjunto de terminales y la cadena vacia `lambda`($\lambda$):
```c
; = { } ( ) return read id write ,  function int float boolean void string intVal realVal stringVal true false - + == ! && != do while if let  *= lambda 
```

A continuacion se definen las producciones de la gramatica:
```c
P -> Decl P  
P -> DecFunc P
P -> lambda

Decl -> if ( Expr ) Sent
Decl -> let TipoDecl id ;
Decl -> do WhileBody
Decl -> Sent 

TipoDecl -> Tipo
TipoDecl -> lambda

WhileBody -> { FuncBody } while ( Expr ) ;

Expr -> ExpRel Expr2

Expr2 ->  && ExpRel Expr2   
Expr2 -> lambda             

ExpRel -> AritExp ExpRel2   

ExpRel2 -> == AritExp ExpRel2   
ExpRel2 -> != AritExp ExpRel2   
ExpRel2 -> lambda               

AritExp -> Term AritExp2        

AritExp2 -> + Term AritExp2     
AritExp2 -> - Term AritExp2     
AritExp2 -> lambda              

Term -> ! Term3                 
Term -> + Term2                 
Term -> - Term2                 
Term -> Term2                   

Term3 -> true       
Term3 -> false      
Term3 -> Term2      

Term2 -> intVal     
Term2 -> realVal    
Term2 -> id FactorId    
Term2 -> stringVal      
Term2 -> ( Expr )       

FactorId -> ( ParamList )       
FactorId -> lambda      

DecFunc -> function TipoFunc id ( FuncParams ) { FuncBody } 

TipoFunc -> Tipo        
TipoFunc -> void        

FuncParams -> Tipo id FuncParams2       
FuncParams -> void              

FuncParams2 -> , Tipo id FuncParams2        
FuncParams2 -> lambda               

Tipo -> int             
Tipo -> float           
Tipo -> boolean         
Tipo -> string          

FuncBody -> Decl FuncBody       
FuncBody -> lambda              

ParamList -> Expr ParamList2        
ParamList -> lambda                 

ParamList2 -> , Expr ParamList2         
ParamList2 -> lambda                

Sent -> id Sent2                
Sent -> write Expr ;            
Sent -> read id ;               
Sent -> return ReturnExp ;      

Sent2 -> = Expr ;           
Sent2 -> *= Expr ;          
Sent2 -> ( ParamList ) ;    

ReturnExp -> Expr       
ReturnExp -> lambda     
```

Siendo P el **axioma** de la gramatica, es decir, el punto de partida.

### Tablas FIRST & FOLLOW
El conjunto **FIRST** es el conjunto de terminales iniciales de una produccion, es decir, es el conjunto de todos los terminales que una produccion $A$ tiene al comienzo(puede incluir a lambda).

El conjunto **FOLLOW** es el conjunto de terminales que pueden venir despues de la produccion $A$. Este conjunto no puede contener lambdas. El axioma ademas contiene al simbolo de fin $\$$.


Se detalla la tabla de **FIRST** & **FOLLOW** para todo *No-Terminal* de la gramática con el objetivo de analizar **si es una gramática adecuada** para un **analizador sintáctico descendente recursivo LL(1).**

|    **NT**     |                   **FIRST**                   |                **FOLLOW**                |
| :-----------: | :-------------------------------------------: | :--------------------------------------: |
|      *P*      |   if, let, do,  id, write, read, return, λ    |                    $                     |
|    *Decl*     |     if, let, do, id, write, read, return      |  if, let, do,  id, write, read, return   |
|  *TipoDecl*   |        int, float, boolean, string, λ         |                    id                    |
|  *WhileBody*  |                       {                       |  if, let, do,  id, write, read, return   |
|     *Exp*     |  +, -, !, intVal, realVal, id, stringVal, (   |                  ), `,`                  |
|    *Exp2*     |                     &&, λ                     |                  ), `,`                  |
|   *ExpRel*    |  +, -, !, intVal, realVal, id, stringVal, (   |                &&, ), `,`                |
|   *ExpRel2*   |                   ==, !=, λ                   |                &&, ), `,`                |
|   *AritExp*   |  +, -, !, intVal, realVal, id, stringVal, (   |            ==, !=,&&, ), `,`             |
|  *AritExp2*   |                    +, -, λ                    |            ==, !=,&&, ), `,`             |
|    *Term*     |  +, -, !, intVal, realVal, id, stringVal, (   |         +, -, ==, !=,&&, ), `,`          |
|    *Term3*    |  true, false, intVal,realVal,id,stringVal,(   |         +, -, ==, !=,&&, ), `,`          |
|    *Term2*    |       intVal, realVal, id, stringVal, (       |         +, -, ==, !=,&&, ), `,`          |
|  *FactorId*   |                     (, λ                      |         +, -, ==, !=,&&, ), `,`          |
|   *DecFunc*   |                   function                    | if, let, do,  id, write, read, return, $ |
|  *TipoFunc*   |       void, int, float, boolean, string       |                    id                    |
| *FuncParams*  |       void, int, float, boolean, string       |                    )                     |
| *FuncParams2* |                    `,`, λ                     |                    )                     |
|    *Tipo*     |          int, float, boolean, string          |                    id                    |
|  *FuncBody*   |    if, let, do, id, write, read, return, λ    |                    }                     |
|  *ParamList*  | +, -, !, intVal, realVal, id, stringVal, (, λ |                    )                     |
| *ParamList2*  |                    `,`, λ                     |                    )                     |
|    *Sent*     |            id, write, read, return            |  if, let, do,  id, write, read, return   |
|    *Sent2*    |                   *=, =, (                    |  if, let, do,  id, write, read, return   |
|  *ReturnExp*  | +, -, !, intVal, realVal, id, stringVal, (, λ |                    ;                     |

Analizando la tabla, podemos observar que las producciones que tienen $\lambda$ en el conjunto **FIRST** tienen intersección nula con el conjunto **FOLLOW** de la misma producción.
Podemos concluir que la gramática es adecuada para LL(1) RD.

## Analizador Semantico
