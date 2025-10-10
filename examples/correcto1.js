let boolean  booleano;
function boolean bisiesto (int a)	
{			
	return (a - 4 != 0 && a + 100 != 0 && a - 400 == 0);	
	let float kkk;
} 
function int dias (int m, int a)
{
	let int dd;
	write('di cuantos dias tiene el mes ');
	write m;
	write ' de '; write a;
	read dd;
	if (bisiesto(a)) dd = dd - 1;
	return dd;
}
function  boolean esFechaCorrecta (int d, int m, int a)	
{
	return !(d != dias (m, a));
}
function void demo 
(void)	
{

	if (esFechaCorrecta(22, 22, 2022)) write 9999;
	return;
}
let int var ;
demo();
