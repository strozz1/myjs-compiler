/******* José Luis Fuertes, septiembre, 2025 *********/
/* El ejemplo incorpora elementos del lenguaje opcionales y elementos que no todos los grupos tienen que implementar */

let string/* variable global cadena */
	s; 	
	
function int FactorialRecursivo (int _n)	/* _n: parámetro formal de la función entera */
{
	if (_n == 0)	return 1;
	return _n * FactorialRecursivo (_n - 1);	/* llamada recursiva */
}

let int uno = 1;	
let int UNO = uno;

function string salto (void)
{
	return '\n';
}

function int FactorialDo (int n)
{
	let int factorial = 0 + uno * 1;	
	do
	{
		factorial *= n--;
	} while (n != 0);		
	return factorial;		
}

function int FactorialWhile (void)
{
	let int factorial = 1;	
	let   int i;			
	while (i < num)			
	{
		factorial *= ++i;	
	}
	return factorial;
}

function int FactorialFor (int n)
{
	let int i;
	let int  factorial= UNO;	/* declaración de variables locales */
	for (i = 1; i <= n; i++)
	{
		factorial *= i;
	}
	return factorial;
}

let int For;
let int Do;
let int While;	

function void imprime (string s, string msg, int f)	/* función que recibe 3 argumentos */
{
	write s;write msg ;write (f);
	write salto();	
	return;	/* finaliza la ejecución de la función (en este caso, se podría omitir) */
}

function string cadena (boolean log_)
{
	if (!log_) {return s;}
	else       {return'Fin';}
}	


s = 'El factorial ';	

write s;
write
 '\nIntroduce un 'número'.';
read 
 num;	/* se lee un número del teclado y se guarda en la variable global num */

switch (num)
{
	case 1:
	case 0: write 'El factorial de '; write num; write' siempre es 1.\n'; break;
	default:
		if (num < 0)
		{
			write ('No existe el factorial de un negativo.\n');
		}
		else
		{
			For = FactorialFor (num);
			While = FactorialWhile ();
			Do = FactorialDo (num);
			imprime (cadena (false), 
					'recursivo es: ', 
					FactorialRecursivo (num));
			imprime (s, 
					'con do-while es: ', 
					Do);
			imprime (s, 
					'con while es: ', 
					While);
			imprime (cadena (false), 
					'con for es: ', 
					For);
		}
}

function boolean  bisiesto(int a)	
{			
	return 
		(a % 4 == 0 && a % 100 != 0 || a % 400 == 0);	
} 

/* OJO:
	- Esto son llaves: {}
	- Esto son corchetes: []
	- Esto son paréntesis: ()
*/

function int dias (int m, int a)
{
	switch (m)
	{
		case 1: case 3: case 5: case 7: case 8: case 10: case 12:
			return 31; break;
		case 4: case 6: case 9: case 11:
			return 30;
		case 2: if (bisiesto (a))  return(29); 
			return(28);
		default: write 'Error: mes incorrecto: '; write m; write salto(); return 0;
	}
} 

function  boolean esFechaCorrecta(int d, int m, int a)	
{
	return m>=1&&m<=12&&d>=1&&d<=dias(m,a);
} 

function void imprimeSuma (int v, int w)	
{
	write v + w;
	write (salto());
} 

function void imprimeReal (float _r_r_)	
{
	write _r_r_;
	write salto();
} 

function void potencia(int z, int dim)	
{
	let int s;	
	for (s=0; s < dim; s++)
	{
		z *= z;
		imprime ('Potencia:', ' ', z);
	}
} 

function void _demo (void)	/* definición de la función _demo, sin argumentos y que no devuelve nada */
{
	let int v1; let int v2; let int v3;
	let int zv ; 
	let float r_1, r_2;
	
	write 'Escribe un número real no muy grande: '; read r_2;
	if (r_2 > 777.777)
	{
		write ('¡Te has pasado!');
	}
	else
	{
		for (r_1= -3.333; r_1 < r_2; r_1+= 3.456789)
		{
			r_2= r_2 / 01.23456789;
			imprimeReal (r_1); 
			write '¿Hola, mundo? ;-) '; 
			imprimeReal (r_2); 
		}
	}

	write'Escribe 'tres' enteros: ';
	read v1; read v2; read v3;
	
	if (v3 == 0) return;
	
	if (!((v1 == v2) && (v1 == v3)))	/* NOT ((v1 igual a v2) AND (v1 distinto de v3))  */
	{
		write 'Escriba su nombre: ';
		let string s;	
		read s;
		if (v2 < v3)	/* si v2<v3, v0=v2; en otro caso v0=1/v3 */
		{
			let int v0=v2; 
		}
		else
		{
			v0= 1 / v3;
		}
		write s;
	}
	s = 'El primer valor era ';
	if (v1 != 0)
	{
		write (s); write v1; write '.\n';
	}
		write s; imprimeSuma (uno, -UNO); write ('.\n');	

	potencia (v0, 4);
	let int i;

	potencia (zv, 5);
	imprimeSuma (i, num);
	imprime ('', cadena(true), 666);
}

_demo();
/* esto constituye la llamada a una función sin argumentos.  */
