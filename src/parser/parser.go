package parser

import (
	"bufio"
	"compiler-pdl/src/errors"
	"compiler-pdl/src/lexer"
	"compiler-pdl/src/token"
	"fmt"
)

var DEBUG bool

type Parser struct {
	parserExec ParserExec
}

func NewParser(lexer *lexer.Lexer) Parser {
	if DEBUG {
		fmt.Println("DEBUG: Initialized Parser")
		fmt.Println("DEBUG: Reading first Token from Lexer")
	}
	lookahead, ok := lexer.Lexical()
	parserExec := ParserExec{
		lexer,
		lookahead,
		[]int{},
	}
	if !ok {
		//todo
		return Parser{}
	}
	return Parser{
		parserExec,
	}
}
func (p *Parser) Parse() bool {
	r := p.parserExec.P(Attr{})
	p.parserExec.lexer.STManager.DestroyScope()
	return (r.tipo == OK)
}

type ParserExec struct {
	lexer     *lexer.Lexer
	lookahead token.Token
	list      []int
}

type Attr struct {
	tipo          TypeExp
	idPos         int
	posActual     int // para params
	numParam      int
	declParamList []string
	paramActual   string
	funcBody      bool
	returnType    TypeExp
}

func error() Attr {
	return Attr{tipo: ERROR}
}

type TypeExp int

const (
	VOID TypeExp = iota
	INT
	FLOAT
	STRING
	FUNCTION
	BOOLEAN
	OK
	ERROR
)

func (t TypeExp) String() string {
	var a string
	switch t {
	case VOID:
		a = "void"
	case INT:
		a = "int"
	case FLOAT:
		a = "float"
	case STRING:
		a = "string"
	case FUNCTION:
		a = "function"
	case BOOLEAN:
		a = "boolean"
	}
	return a
}
func from(s string) TypeExp {
	var t TypeExp
	switch s {
	case "void":
		t = VOID
	case "int":
		t = INT
	case "float":
		t = FLOAT
	case "string":
		t = STRING
	case "function":
		t = FUNCTION
	case "boolean":
		t = BOOLEAN
	}
	return t
}

// Function that saves wich rule has been applied
func (p *ParserExec) rule(i int) {
	p.list = append(p.list, i)
}

func (p *Parser) Write(w *bufio.Writer) {
	if DEBUG {
		fmt.Println("DEBUG: Writting parse to output")
	}

	fmt.Fprintf(w, "D ")
	for _, i := range p.parserExec.list {
		fmt.Fprintf(w, "%d ", i)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func (p *ParserExec) getToken() bool {
	tk, ok := p.lexer.Lexical()
	if ok {
		p.lookahead = tk

	}
	return ok
}

// match token and ask for next token to lexer, saving it in p.lookahead
// returns true if match is successful
func (p *ParserExec) match(tk token.TokenKind, attr any) bool {
	if p.lookahead.Kind == tk {
		if attr != nil {
			if p.lookahead.Attr == attr {
				p.getToken()
				return true
			}
			return false
		}
		p.getToken()
		return true
	}
	return false
}

// Axiom
func (p *ParserExec) P(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.LET, token.ID, token.IF, token.DO, token.READ, token.WRITE, token.RETURN:
		p.rule(1)

		res := p.Decl(attr)
		if res.tipo == ERROR {
			return error()
		}
	case token.FUNCTION:
		p.rule(2)
		res := p.DecFunc(attr)
		if res.tipo == ERROR {
			return error()
		}
		p.lexer.STManager.DestroyScope()
	case token.EOF:
		p.rule(3)
		return Attr{tipo: OK}
	default:
		if(attr.tipo!=ERROR){
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 3")
		}
		return error()
	}
	t := p.lookahead.Kind
	if !(t == token.FUNCTION || t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN || p.lexer.EOF) {
		errors.SintacticalError(errors.S_EXPECTED_SENT, nil)
		//os.Exit(1);
		return error()
	}
	res := p.P(attr)
	if res.tipo == ERROR {
		return error()
	}
	return res
}

func (p *ParserExec) Decl(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.IF:
		p.rule(4)
		p.match(token.IF, nil)
		if p.lookahead.Kind != token.ABRIR_PAR {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 4")
			return error()
		}
		p.match(token.ABRIR_PAR, nil)
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 4")
			return error()
		}
		exp := p.Expr(attr)
		if exp.tipo == ERROR {
			return error()
		}
		if exp.tipo != BOOLEAN {
			errors.SemanticalError(errors.SS_IF_COND, exp.tipo.String())
			return error()
		}
		if p.lookahead.Kind != token.CERRAR_PAR {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 4")
			return error()
		}
		p.match(token.CERRAR_PAR, nil)
		t = p.lookahead.Kind
		if !(t == token.ID || t == token.WRITE || t == token.READ || t == token.RETURN) {
			errors.SintacticalError(errors.S_EXPECTED_CERRAR_PAR, nil)
			return error()
		}
		return p.Sent(attr)
	case token.LET:
		p.lexer.DeclarationZone(true)
		p.rule(5)
		p.match(token.LET, nil)
		t := p.lookahead.Kind
		if !(t == token.ID || t == token.INT || t == token.FLOAT || t == token.BOOLEAN ||
			t == token.STRING) {
			errors.SintacticalError(errors.S_TYPE, nil) // 5")
			return error()
		}
		tipoDecl := p.TipoDecl().tipo
		if(tipoDecl==ERROR){
			return error()
		}
		if p.lookahead.Kind != token.ID {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 5")
			return error()
		}
		pos := p.lookahead.Attr.(int)
		entry, ok := p.lexer.STManager.GetEntry(pos)
		if !ok {
			//todo
			errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		p.lexer.STManager.SetEntryType(entry, tipoDecl.String())
		p.match(token.ID, nil)
		p.lexer.DeclarationZone(false)
		if p.lookahead.Kind != token.PUNTOYCOMA {
			errors.SintacticalError(errors.S_EXPECTED_SEMICOLON, nil) // 5")
			return error()
		}
		p.match(token.PUNTOYCOMA, nil)
	case token.DO:
		p.rule(6)
		p.match(token.DO, nil)
		if p.lookahead.Kind != token.ABRIR_CORCH {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 6")
			return error()
		}
		res := p.WhileBody(attr)
		return res
	default:
		t := p.lookahead.Kind
		if t == token.ID || t == token.WRITE || t == token.READ || t == token.RETURN {
			p.rule(7)
			return p.Sent(attr)
		} else {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 7")
			return error()
		}
	}
	return Attr{tipo: OK}
}

func (p *ParserExec) TipoDecl() Attr {
	t := p.lookahead.Kind
	if t == token.INT || t == token.FLOAT || t == token.BOOLEAN || t == token.STRING {
		//if FIRST(Tipo)
		p.rule(8)
		return p.Tipo()
	} else if t != token.ID {
		p.rule(9)
		return Attr{tipo: INT}
	} else {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 9")
		return error()
	}
}

func (p *ParserExec) WhileBody(attr Attr) Attr {
	p.rule(10)
	if p.lookahead.Kind != token.ABRIR_CORCH {
		errors.SintacticalError(errors.S_EXPECTED_WHILE_CORCH, nil) // 10")
		return error()
	}
	p.match(token.ABRIR_CORCH, nil)
	t := p.lookahead.Kind
	if !(t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN ||
		t == token.CERRAR_CORCH) {
		errors.SintacticalError(errors.S_EXPECTED_SENT, nil) // 10")
		return error()
	}
	bodyAttr := p.FuncBody(attr)
	if(bodyAttr.tipo==ERROR){
		return bodyAttr
	}
	if p.lookahead.Kind != token.CERRAR_CORCH {
		errors.SintacticalError(errors.S_EXPECTED_CERRAR_CORCH, nil) // 10")
		return error()
	}
	p.match(token.CERRAR_CORCH, nil)
	if p.lookahead.Kind != token.WHILE {
		errors.SintacticalError(errors.S_MISSING_WHILE, nil) // 10")
		return error()
	}
	p.match(token.WHILE, nil)
	if p.lookahead.Kind != token.ABRIR_PAR {
		errors.SintacticalError(errors.S_EXPECTED_ABRIR_PAR, nil) // 10")
		return error()
	}
	p.match(token.ABRIR_PAR, nil)
	t = p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 10")
		return error()
	}
	condAttr := p.Expr(attr)
	if(condAttr.tipo==ERROR){
		return condAttr;
	}
	if p.lookahead.Kind != token.CERRAR_PAR {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 10")
		return error()
	}
	p.match(token.CERRAR_PAR, nil)
	if p.lookahead.Kind != token.PUNTOYCOMA {
		errors.SintacticalError(errors.S_EXPECTED_SEMICOLON, nil) // 10")
		return error()
	}
	p.match(token.PUNTOYCOMA, nil)
	if bodyAttr.tipo == ERROR {
		return error()
	}
	if condAttr.tipo == ERROR {
		return error()
	}
	if condAttr.tipo != BOOLEAN {
		errors.SemanticalError(errors.SS_EXPECTED_WHILE_COND, condAttr.tipo.String())
		return error()
	}
	return Attr{tipo: OK}
}

// Expr -> ExpRel {if(ExpRel!=error)Expr2.tipo:=ExpRel.tipo else Expr.error} Expr2 {Expr.tipo=Expr2.tipo}
func (p *ParserExec) Expr(attr Attr) Attr {
	p.rule(11)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 11")
		return error()
	}
	relRes := p.ExpRel(attr)
	if relRes.tipo == ERROR {
		return relRes
	}

	t = p.lookahead.Kind
	if !(t == token.ARITM || t == token.LOGICO || t == token.CERRAR_PAR ||
		t == token.PUNTOYCOMA || t == token.COMA) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 11b")
		return error()
	}

	a := p.Expr2(relRes)
	return a
}

func (p *ParserExec) Expr2(attr Attr) Attr {
	t := p.lookahead.Kind
	if t == token.LOGICO && p.lookahead.Attr == token.LOG_AND {
		p.rule(12)
		p.match(token.LOGICO, token.LOG_AND)

		t = p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP_LOG, nil) // 11")
			return error()
		}

		resRel := p.ExpRel(attr)
		if resRel.tipo == ERROR {
			return resRel
		} else if !(resRel.tipo == BOOLEAN && attr.tipo == BOOLEAN) {
			errors.SemanticalError(errors.SS_EXPECTED_BOOLEANS, fmt.Sprintf("%s && %s", resRel.tipo.String(), attr.tipo.String()))
			return error()
		}

		t = p.lookahead.Kind
		if !(t == token.ARITM || t == token.LOGICO || t == token.CERRAR_PAR ||
			t == token.PUNTOYCOMA || t == token.COMA) {
			errors.SintacticalError(errors.S_EXPECTED_EXP_LOG, nil) // 11b")
			return error()
		}
		return p.Expr2(resRel)
	} else if t == token.ARITM || t == token.LOGICO || t == token.CERRAR_PAR ||
		t == token.PUNTOYCOMA || t == token.COMA {
		p.rule(13) //lambda
		return attr
	} else {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 12/13")
		return error()
	}
}

func (p *ParserExec) ExpRel(attr Attr) Attr {
	p.rule(14)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 14")
		return error()
	}

	resArit := p.AritExp(attr)
	if resArit.tipo == ERROR {
		return resArit
	}

	t = p.lookahead.Kind
	if !(t == token.RELAC || t == token.LOGICO || t == token.ARITM || t == token.COMA ||
		t == token.PUNTOYCOMA || t == token.CERRAR_PAR) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 14")
		return error()
	}
	return p.ExpRel2(resArit)
}
func (p *ParserExec) ExpRel2(attr Attr) Attr {

	switch p.lookahead.Kind {
	case token.RELAC:
		switch p.lookahead.Attr {
		case token.REL_EQ:
			p.rule(15)
			p.match(token.RELAC, token.REL_EQ)
		case token.REL_NOTEQ:
			p.rule(16)
			p.match(token.RELAC, token.REL_NOTEQ)
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 15/16")
			return error()
		}
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 15/16")
			return error()
		}

		arit := p.AritExp(attr)
		if arit.tipo == ERROR {
			return error()
		}
		if !(arit.tipo.String() == attr.tipo.String()) {
			errors.SemanticalError(errors.SS_RELATIONAL_TYPES, fmt.Sprintf("'%s' y '%s'", attr.tipo.String(), arit.tipo.String()))
			return error()
		}

		t = p.lookahead.Kind
		if !(t == token.RELAC || t == token.LOGICO || t == token.ARITM || t == token.COMA ||
			t == token.PUNTOYCOMA || t == token.CERRAR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 15/16b")
			return error()
		}
		arit.tipo = BOOLEAN
		return p.ExpRel2(arit)

	case token.ARITM, token.COMA, token.PUNTOYCOMA, token.LOGICO, token.CERRAR_PAR:
		p.rule(17)
		return attr
	default:
		if(attr.tipo!=ERROR){
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 17")
		}
		return error()
	}
}

func (p *ParserExec) AritExp(attr Attr) Attr {
	p.rule(18)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 18/a")
		return error()
	}
	termAttr := p.Term(attr)
	if termAttr.tipo == ERROR {
		return termAttr
	}

	t = p.lookahead.Kind
	if !(t == token.LOGICO || t == token.ARITM || t == token.COMA ||
		t == token.PUNTOYCOMA || t == token.CERRAR_PAR || t == token.RELAC) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 18b")
		return error()
	}
	return p.AritExp2(termAttr)
}

func (p *ParserExec) AritExp2(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.ARITM:
		switch p.lookahead.Attr {
		case token.ARIT_PLUS:
			p.match(token.ARITM, token.ARIT_PLUS)
			p.rule(19)
		case token.ARIT_MINUS:
			p.rule(20)
			p.match(token.ARITM, token.ARIT_MINUS)
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 20")
			return error()
		}
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 19")
			return error()
		}
		termRes := p.Term(attr)
		if termRes.tipo == ERROR {
			return error()
		}
		if !((attr.tipo == FLOAT && termRes.tipo == FLOAT) || (attr.tipo == INT && termRes.tipo == INT)) {
			errors.SemanticalError(errors.SS_INVALID_ARIT_TYPES, fmt.Sprintf("'%s' y '%s'", attr.tipo.String(), termRes.tipo.String()))
			return error()
		}

		t = p.lookahead.Kind
		if !(t == token.LOGICO || t == token.ARITM || t == token.COMA ||
			t == token.PUNTOYCOMA || t == token.CERRAR_PAR || t == token.RELAC) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 19/20")
			return error()
		}
		return p.AritExp2(termRes)
	case token.RELAC, token.LOGICO, token.CERRAR_PAR, token.PUNTOYCOMA, token.COMA:
		p.rule(21) //lambda
		return attr
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil)
		return error()
	}
}

func (p *ParserExec) Term(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.LOGICO:
		if p.lookahead.Attr != token.LOG_NEG {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 22")
			return error()
		}
		p.rule(22)
		p.match(token.LOGICO, token.LOG_NEG)

		t := p.lookahead.Kind
		if !(t == token.TRUE || t == token.FALSE || t == token.INT_LITERAL ||
			t == token.REAL_LITERAL || t == token.STRING_LITERAL ||
			t == token.ABRIR_PAR || t == token.ID) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 22")
			return error()
		}
		res := p.Term3(attr)
		if res.tipo == ERROR {
			return error()
		}
		if res.tipo != BOOLEAN {
			errors.SemanticalError(errors.SS_NEGATION_EXPECTED_BOOL, res.tipo.String())
			return error()
		}
		return res
	case token.ARITM:
		switch p.lookahead.Attr {
		case token.ARIT_PLUS:
			p.rule(23)
			p.match(token.ARITM, token.ARIT_PLUS)
		case token.ARIT_MINUS:
			p.rule(24)
			p.match(token.ARITM, token.ARIT_MINUS)
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 23/24")
			return error()
		}

		switch p.lookahead.Kind {
		case token.INT_LITERAL, token.REAL_LITERAL, token.ID, token.STRING_LITERAL,
			token.ABRIR_PAR:
			term2 := p.Term2(attr)
			if term2.tipo == ERROR {
				return error()
			}
			if !(term2.tipo == INT || term2.tipo == FLOAT) {
				errors.SemanticalError(errors.SS_INVALID_SIGN_TYPE, term2.tipo.String())
				return error()
			}
			return term2
		default:
			
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 23/24")
			return error()
		}
	case token.INT_LITERAL, token.REAL_LITERAL, token.ID, token.STRING_LITERAL,
		token.ABRIR_PAR:
		p.rule(25)
		return p.Term2(attr)
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 23/24")
		return error()
	}
}

func (p *ParserExec) Term3(attr Attr) Attr {
	var t TypeExp
	switch p.lookahead.Kind {
	case token.TRUE:
		p.rule(26)
		p.match(token.TRUE, nil)
		t = BOOLEAN
	case token.FALSE:
		p.rule(27)
		p.match(token.FALSE, nil)
		t = BOOLEAN
	case token.INT_LITERAL, token.REAL_LITERAL, token.ID, token.STRING_LITERAL,
		token.ABRIR_PAR:
		p.rule(28)
		return p.Term2(attr)

	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 28")
		t = ERROR
	}
	return Attr{tipo: t}
}

func (p *ParserExec) Term2(attr Attr) Attr {
	var tt TypeExp
	switch p.lookahead.Kind {
	case token.INT_LITERAL:
		p.rule(29)
		p.match(token.INT_LITERAL, nil)
		tt = INT
	case token.REAL_LITERAL:
		p.rule(30)
		p.match(token.REAL_LITERAL, nil)
		tt = FLOAT
	case token.ID:
		p.rule(31)
		pos := p.lookahead.Attr.(int)
		attr.idPos = pos
		p.match(token.ID, nil)
		t := p.lookahead.Kind
		if !(t == token.ABRIR_PAR || t == token.CERRAR_PAR || t == token.ARITM ||
			t == token.RELAC || t == token.LOGICO || t == token.COMA ||
			t == token.PUNTOYCOMA) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 31")
			return error()
		}
		res := p.FactorId(attr)
		if res.tipo == ERROR {
			return res
		}
		tt = res.tipo
	case token.STRING_LITERAL:
		p.rule(32)
		p.match(token.STRING_LITERAL, nil)
		tt = STRING
	case token.ABRIR_PAR:
		p.rule(33)
		p.match(token.ABRIR_PAR, nil)

		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 33")
			return error()
		}

		res := p.Expr(attr)
		if res.tipo == ERROR {
			return error()
		}
		tt = res.tipo

		if p.lookahead.Kind != token.CERRAR_PAR {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 33")
			return error()
		}
		p.match(token.CERRAR_PAR, nil)
	}
	return Attr{tipo: tt}
}

func (p *ParserExec) FactorId(attr Attr) Attr {
	entry, ok := p.lexer.STManager.GetEntry(attr.idPos)
	if !ok {
		errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
		return error()
	}
	switch p.lookahead.Kind {
	case token.ABRIR_PAR:
		p.rule(34)
		p.match(token.ABRIR_PAR, nil)
		if entry.GetType().String() != FUNCTION.String() {
			errors.SemanticalError(errors.SS_FUNC_NOT_DECL, entry.Lexeme)
			return error()
		}
		switch p.lookahead.Kind {
		case token.ARITM, token.LOGICO, token.CERRAR_PAR, token.ID, token.INT_LITERAL,
			token.REAL_LITERAL, token.STRING_LITERAL:
			if p.lookahead.Kind == token.LOGICO {
				if p.lookahead.Attr != token.LOG_NEG {
					errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 34a")
					return error()
				}
			}
			num := entry.GetAttribute("numParam").Value().(int)
			attr.posActual = 0
			attr.numParam = num
			res := p.ParamList(attr)
			if res.tipo == ERROR {
				return res
			}
			if p.lookahead.Kind != token.CERRAR_PAR {
				errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 34c")
				return error()
			}
			p.match(token.CERRAR_PAR, nil)
			if res.tipo == OK {
				ret := entry.GetAttribute("tipoRetorno").Value().(string)
				res.tipo = from(ret)
				return res
			} else {
				return Attr{tipo: ERROR}
			}
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 34b")
			return error()
		}
	case token.LOGICO, token.ARITM, token.RELAC, token.COMA, token.PUNTOYCOMA,
		token.CERRAR_PAR:
		p.rule(35)
		t := from(entry.GetType().String())
		attr.tipo = t
		return attr
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 35")
		return error()
	}
}

func (p *ParserExec) DecFunc(attr Attr) Attr {
	p.rule(36)
	if p.lookahead.Kind != token.FUNCTION {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}
	p.lexer.DeclarationZone(true)
	p.match(token.FUNCTION, nil)
	t := p.lookahead.Kind
	if !(t == token.STRING || t == token.VOID || t == token.INT ||
		t == token.FLOAT || t == token.BOOLEAN) {
		errors.SintacticalError(errors.S_EXPECTED_FUNCTYPE, nil) // 36")
		return error()
	}
	tipo := p.TipoFunc()
	if tipo.tipo == ERROR {
		return error()
	}
	if p.lookahead.Kind != token.ID {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}
	i, _ := p.lookahead.Attr.(int)
	e, ok := p.lexer.STManager.GetEntry(i)
	if !ok {
		errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
		return error()
	}
	funcName := p.lookahead.Lexeme
	p.lexer.STManager.SetEntryType(e, "function")
	e.SetAttributeValue("tipoRetorno", tipo.tipo.String())
	attr.idPos = i
	attr.returnType = tipo.tipo
	p.match(token.ID, nil)
	if p.lookahead.Kind != token.ABRIR_PAR {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}
	newScope := fmt.Sprintf("Function %s", funcName)
	p.lexer.STManager.NewScope(newScope)
	p.match(token.ABRIR_PAR, nil)
	t = p.lookahead.Kind
	if !(t == token.STRING || t == token.VOID || t == token.INT ||
		t == token.FLOAT || t == token.BOOLEAN) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}
	attr.numParam = 0
	attr.declParamList = []string{}
	params := p.FuncParams(attr)
	if p.lookahead.Kind != token.CERRAR_PAR {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}
	if params.tipo == ERROR {
		return error()
	}
	p.match(token.CERRAR_PAR, nil)
	if p.lookahead.Kind != token.ABRIR_CORCH {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}
	p.lexer.DeclarationZone(false)
	p.match(token.ABRIR_CORCH, nil)
	t = p.lookahead.Kind
	if !(t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN ||
		t == token.CERRAR_CORCH) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}

	params.funcBody = true
	bodyRes := p.FuncBody(params)
	bodyRes.funcBody = false
	if p.lookahead.Kind != token.CERRAR_CORCH {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 36")
		return error()
	}
	p.match(token.CERRAR_CORCH, nil)

	return bodyRes
}

func (p *ParserExec) TipoFunc() Attr {
	switch p.lookahead.Kind {
	case token.VOID:
		p.rule(38)
		p.match(token.VOID, nil)
		return Attr{tipo: VOID}
	case token.STRING, token.INT, token.FLOAT, token.BOOLEAN:
		p.rule(37)
		return p.Tipo()
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 37/38")
		return error()
	}
}

func (p *ParserExec) FuncParams(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.INT, token.FLOAT, token.BOOLEAN, token.STRING: //FIRST(Tipo)
		p.rule(39)
		tipo := p.Tipo()
		if tipo.tipo == ERROR {
			return error()
		}
		attr.declParamList = append(attr.declParamList, tipo.tipo.String())
		attr.numParam++
		if p.lookahead.Kind != token.ID {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 39")
			return error()
		}
		i, _ := p.lookahead.Attr.(int)
		e, ok := p.lexer.STManager.GetEntry(i)
		if !ok {
			errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		p.lexer.STManager.SetEntryType(e, tipo.tipo.String())
		p.match(token.ID, nil)
		if p.lookahead.Kind == token.COMA || p.lookahead.Kind == token.CERRAR_PAR {
			//first FuncParams2
			return p.FuncParams2(attr)
		} else {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 39")
			return error()
		}
	case token.VOID:
		p.match(token.VOID, nil)
		p.rule(40)
		e, ok := p.lexer.STManager.GetGlobalEntry(attr.idPos)
		if !ok {
			errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		e.SetAttributeValue("numParam", attr.numParam) //0
		//e.SetAttributeValue("tipoParam1", attr.tipoParam)
		attr.tipo = OK
		return attr

	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 40")
		return error()
	}
}

func (p *ParserExec) FuncParams2(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.COMA:
		p.match(token.COMA, nil)
		p.rule(41)
		switch p.lookahead.Kind {
		case token.INT, token.FLOAT, token.BOOLEAN, token.STRING, token.ID: //FIRST(Tipo)&id
			tipo := p.Tipo()
			if tipo.tipo == ERROR {
				return error()
			}
			attr.declParamList = append(attr.declParamList, tipo.tipo.String())
			attr.numParam++
			if p.lookahead.Kind != token.ID {
				errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 41a")
				return error()
			}
			i, _ := p.lookahead.Attr.(int)
			e, ok := p.lexer.STManager.GetEntry(i)
			if !ok {
				errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
				return error()
			}
			p.lexer.STManager.SetEntryType(e, tipo.tipo.String())
			p.match(token.ID, nil)
			if !(p.lookahead.Kind == token.COMA || p.lookahead.Kind == token.CERRAR_PAR) {
				errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 41b")
				return error()
			}
			return p.FuncParams2(attr)
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 41")
			return error()
		}
	case token.CERRAR_PAR:
		p.rule(42)
		e, ok := p.lexer.STManager.GetGlobalEntry(attr.idPos)
		if !ok {
			errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		e.SetAttributeValue("numParam", attr.numParam)
		for a, b := range attr.declParamList {
			name := fmt.Sprintf("tipoParam%d", a+1)
			p.lexer.STManager.SetEntryAttribute(e, name, b)
		}
		attr.tipo = OK
		return attr

	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 42")
		return error()
	}
}

func (p *ParserExec) Tipo() Attr {
	var t TypeExp
	switch p.lookahead.Kind {
	case token.INT:
		p.match(token.INT, nil)
		p.rule(43)
		t = INT
	case token.FLOAT:
		p.match(token.FLOAT, nil)
		p.rule(44)
		t = FLOAT
	case token.BOOLEAN:
		p.match(token.BOOLEAN, nil)
		p.rule(45)
		t = BOOLEAN
	case token.STRING:
		p.match(token.STRING, nil)
		p.rule(46)
		t = STRING
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil)
		t = ERROR
	}
	return Attr{tipo: t}
}

func (p *ParserExec) FuncBody(attr Attr) Attr {

	switch p.lookahead.Kind {
	case token.IF, token.LET, token.DO, token.ID, token.READ, token.WRITE,
		token.RETURN: //first Decl
		p.rule(47)
		res := p.Decl(attr)
		if res.tipo == ERROR {
			return error()
		}
		t := p.lookahead.Kind
		if !(t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
			t == token.WRITE || t == token.READ || t == token.RETURN ||
			t == token.CERRAR_CORCH) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil)
			return error()
		}
		return p.FuncBody(attr)
	case token.CERRAR_CORCH:
		p.rule(48)
		attr.tipo = OK
		return attr
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil)
		return error()
	}
}

func (p *ParserExec) ParamList(attr Attr) Attr {
	t := p.lookahead.Kind
	if t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR {
		p.rule(49)

		if attr.posActual+1 > attr.numParam {
			if attr.numParam == 0 {
				errors.SemanticalError(errors.SS_NUM_PARAMS_INV, "no se esperaban parametros")
			} else {
				errors.SemanticalError(errors.SS_NUM_PARAMS_INV, fmt.Sprintf("Se esperaban %d parametros", attr.numParam))
			}
			return error()
		}

		entry, ok := p.lexer.STManager.GetEntry(attr.idPos)
		if !ok {
			errors.NewError(errors.SEMANTICAL, errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}

		expected := entry.GetAttribute("tipoParam1")
		if expected == nil {
			errors.SemanticalError(errors.SS_NUM_PARAMS_INV, "se esperaba un parametro")
			return error()
		}
		exp := expected.Value().(string)
		res := p.Expr(attr)

		if res.tipo == ERROR {
			return error()
		} else if exp != res.tipo.String() {
			errors.SemanticalError(errors.SS_INVALID_EXP_TYPE, exp)
			return error()
		}
		attr.posActual++

		t = p.lookahead.Kind
		if !(t == token.COMA || t == token.CERRAR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 49")
			return error()
		}
		return p.ParamList2(attr)

	} else if p.lookahead.Kind == token.CERRAR_PAR {
		p.rule(50)
		if attr.posActual < attr.numParam {
			errors.SemanticalError(errors.SS_NUM_PARAMS_INV, "No se esperaban parametros")
			return error()
		}
	} else {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 50")
		return error()
	}
	attr.tipo = OK
	return attr
}

func (p *ParserExec) ParamList2(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.COMA:
		p.match(token.COMA, nil)
		p.rule(51)

		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 51")
			return error()
		}

		if attr.posActual+1 > attr.numParam {
			errors.SemanticalError(errors.SS_NUM_PARAMS_INV, fmt.Sprintf(" Se esperaban %d parametros", attr.numParam))
			return error()
		}
		entry, ok := p.lexer.STManager.GetEntry(attr.idPos)
		if !ok {
			errors.NewError(errors.SEMANTICAL, errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}

		expected := entry.GetAttribute(fmt.Sprintf("tipoParam%d", attr.posActual+1))
		if expected == nil {
			errors.SemanticalError(errors.SS_NUM_PARAMS_INV, "se esperaba un parametro")
			return error()
		}
		exp := expected.Value().(string)
		res := p.Expr(attr)

		if res.tipo == ERROR {
			return error()
		} else if exp != res.tipo.String() {
			errors.SemanticalError(errors.SS_INVALID_EXP_TYPE, expected)
			return error()
		}
		attr.posActual++

		t = p.lookahead.Kind
		if !(t == token.COMA || t == token.CERRAR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 51")
			return error()
		}
		return p.ParamList2(attr)
	case token.CERRAR_PAR:
		p.rule(52)

		if attr.posActual < attr.numParam {
			errors.SemanticalError(errors.SS_NUM_PARAMS_INV, fmt.Sprintf(" Se esperaban %d parametros", attr.numParam))
			return error()
		}
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 52")
		return error()
	}
	attr.tipo = OK
	return attr
}

func (p *ParserExec) Sent(attr Attr) Attr {
	switch p.lookahead.Kind {
	case token.ID:
		i, _ := p.lookahead.Attr.(int)
		a, ok := p.lexer.STManager.GetEntry(i)
		if !ok {
			errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		attr.idPos = i
		attr.tipo = from(a.GetType().String())
		p.match(token.ID, nil)
		p.rule(53)
		switch p.lookahead.Kind {
		case token.ASIG, token.ABRIR_PAR:
			break
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 53")
			return error()
		}
		return p.Sent2(attr)

	case token.WRITE:
		p.match(token.WRITE, nil)
		p.rule(54)
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 54")
			return error()
		}
		exp := p.Expr(attr)
		if p.lookahead.Kind != token.PUNTOYCOMA {
			errors.SintacticalError(errors.S_EXPECTED_SEMICOLON, nil)
			return error()
		}
		p.match(token.PUNTOYCOMA, nil)
		if exp.tipo == ERROR {
			return error()
		}
		attr.tipo = OK
		return attr

	case token.READ:
		p.rule(55)
		if !p.match(token.READ, nil) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 55")
			return error()
		}
		if p.lookahead.Kind != token.ID {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 55")
			return error()
		}
		i, _ := p.lookahead.Attr.(int)
		if !p.lexer.STManager.EntryExists(i) {
			errors.SemanticalError(errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		//TODO: tipo de lectura?
		p.match(token.ID, nil)
		if !p.match(token.PUNTOYCOMA, nil) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 55")
			return error()
		}
		attr.tipo = OK
		return attr

	case token.RETURN:
		p.rule(56)
		p.match(token.RETURN, nil)
		if !(attr.funcBody) {
			errors.SemanticalError(errors.SS_RETURN_OUTSIDE, nil)
			return error()
		}
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.STRING_LITERAL || t == token.REAL_LITERAL ||
			t == token.ID || t == token.ABRIR_PAR || t == token.PUNTOYCOMA) {
			errors.SintacticalError(errors.S_EXPECTED_RET_EXP, nil) // 55")
			return error()
		}
		res := p.ReturnExp()
		if !p.match(token.PUNTOYCOMA, nil) {
			errors.SintacticalError(errors.S_EXPECTED_SEMICOLON, nil) // 55")
			return error()
		}
		if res.tipo != attr.returnType {
			errors.SemanticalError(errors.SS_INVALID_RETURN, fmt.Sprintf("se esperaba %s, se obtuvo %s", attr.returnType.String(), res.tipo.String()))
			return error()
		}
		return res

	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 53-56")
		return error()
	}

}
func (p *ParserExec) Sent2(attr Attr) Attr {

	switch p.lookahead.Kind {
	case token.ASIG:
		switch p.lookahead.Attr {
		case token.ASIG_SIMPLE:
			p.match(token.ASIG, token.ASIG_SIMPLE)
			p.rule(57)
		case token.ASIG_MULT:
			p.match(token.ASIG, token.ASIG_MULT)
			p.rule(58)
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 57/58")
			return error()
		}

		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 57/58")
			return error()
		}

		var res = p.Expr(attr)
		if res.tipo == ERROR {
			return res
		}
		if p.lookahead.Kind != token.PUNTOYCOMA {
			errors.SintacticalError(errors.S_EXPECTED_SEMICOLON, nil) // 57")
			return error()
		}
		p.match(token.PUNTOYCOMA, nil)

		entry, ok := p.lexer.STManager.GetEntry(attr.idPos)
		if !ok {
			errors.NewError(errors.SEMANTICAL, errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		at := entry.GetType()
		if at.String() == res.tipo.String() {
			attr.tipo = OK
			return attr
		} else {
			errors.SemanticalError(errors.SS_INVALID_EXP_TYPE, fmt.Sprintf("%s, se obtuvo %s", at.String(), res.tipo.String()))
			return error()
		}

	case token.ABRIR_PAR:
		p.match(token.ABRIR_PAR, nil)
		p.rule(59)
		switch p.lookahead.Kind {
		case token.ARITM, token.LOGICO, token.CERRAR_PAR, token.ID, token.INT_LITERAL,
			token.REAL_LITERAL, token.STRING_LITERAL:
			if p.lookahead.Kind == token.LOGICO {
				if p.lookahead.Attr != token.LOG_NEG {
					errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 59")
					return error()
				}
			}
		default:
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 59")
			return error()
		}
		if attr.tipo.String() != FUNCTION.String() {
			errors.SemanticalError(errors.SS_EXPECTED_FUNC, attr.tipo.String())
			return error()
		}
		entry, ok := p.lexer.STManager.GetEntry(attr.idPos)
		if !ok {
			errors.NewError(errors.SEMANTICAL, errors.SS_ID_NOT_FOUND, p.lookahead.Lexeme)
			return error()
		}
		numParam := entry.GetAttribute("numParam").Value().(int)
		retType := entry.GetAttribute("tipoRetorno").Value().(string)
		attr.returnType = from(retType)
		attr.numParam = numParam

		var res = p.ParamList(attr)
		if res.tipo == ERROR {
			return res
		}
		if p.lookahead.Kind != token.CERRAR_PAR {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 59")
			return error()
		}
		p.match(token.CERRAR_PAR, nil)

		if p.lookahead.Kind != token.PUNTOYCOMA {
			errors.SintacticalError(errors.S_EXPECTED_SEMICOLON, nil) // 59")
			return error()
		}
		p.match(token.PUNTOYCOMA, nil)
		return res
	default:
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 57-59"
		return error()
	}
}
func (p *ParserExec) ReturnExp() Attr {
	switch p.lookahead.Kind {
	case token.LOGICO:
		if p.lookahead.Attr != token.LOG_NEG {
			errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 60/61")
			return error()
		}
	case token.ARITM, token.INT_LITERAL, token.REAL_LITERAL, token.STRING_LITERAL,
		token.ID, token.ABRIR_PAR:
		break
	case token.PUNTOYCOMA:
		p.rule(61)
		return Attr{tipo: VOID}
	default:
		errors.SintacticalError(errors.S_INVALID_EXP, nil) // 61")
		return error()
	}
	p.rule(60)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.SintacticalError(errors.S_EXPECTED_EXP, nil) // 60")
		return error()
	}
	return p.Expr(Attr{})

}
