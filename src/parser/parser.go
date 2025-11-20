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
func (p *Parser) Parse() {
	p.parserExec.P()
}

type ParserExec struct {
	lexer     *lexer.Lexer
	lookahead token.Token
	list      []int
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
func (p *ParserExec) P() {
	switch p.lookahead.Kind {
	case token.LET, token.ID, token.IF, token.DO, token.READ, token.WRITE, token.RETURN:
		p.rule(1)
		p.Decl()
	case token.FUNCTION:
		p.rule(2)
		p.DecFunc()
	case token.EOF:
		p.rule(3)
		return
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 3")
		return
	}
	t := p.lookahead.Kind
	if !(t == token.FUNCTION || t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN || p.lexer.EOF) {
		fmt.Println("Error 1")
		return
	}
	p.P()
}

func (p *ParserExec) Decl() {
	switch p.lookahead.Kind {
	case token.IF:
		p.rule(4)
		p.match(token.IF, nil)
		if p.lookahead.Kind != token.ABRIR_PAR {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 4")
			return
		}
		p.match(token.ABRIR_PAR, nil)
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 4")
			return
		}
		p.Expr()
		if p.lookahead.Kind != token.CERRAR_PAR {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 4")
			return
		}
		p.match(token.CERRAR_PAR, nil)
		t = p.lookahead.Kind
		if !(t == token.ID || t == token.WRITE || t == token.READ || t == token.RETURN) {
			fmt.Println("Error 4")
			return
		}
		p.Sent()

	case token.LET:
		p.rule(5)
		p.match(token.LET, nil)
		t := p.lookahead.Kind
		if !(t == token.ID || t == token.INT || t == token.FLOAT || t == token.BOOLEAN ||
			t == token.STRING) {
			fmt.Println("Error 5")
			return
		}
		p.TipoDecl()
		if p.lookahead.Kind != token.ID {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 5")
			return
		}
		p.match(token.ID, nil)
		if p.lookahead.Kind != token.PUNTOYCOMA {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 5")
			return
		}
		p.match(token.PUNTOYCOMA, nil)

	case token.DO:
		p.rule(6)
		p.match(token.DO, nil)
		if p.lookahead.Kind != token.ABRIR_CORCH {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 6")
			return
		}
		p.WhileBody()
	default:
		t := p.lookahead.Kind
		if t == token.ID || t == token.WRITE || t == token.READ || t == token.RETURN {
			p.rule(7)
			p.Sent()
		} else {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 7")
			return
		}
	}

}

func (p *ParserExec) TipoDecl() {
	t := p.lookahead.Kind
	if t == token.INT || t == token.FLOAT || t == token.BOOLEAN || t == token.STRING {
		//if FIRST(Tipo)
		p.rule(8)
		p.Tipo()
	} else if t != token.ID {
		p.rule(9)
	} else {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 9")
		return
	}
}

func (p *ParserExec) WhileBody() {
	p.rule(10)
	if p.lookahead.Kind != token.ABRIR_CORCH {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_WHILE_CORCH, nil) // 10")
		return
	}
	p.match(token.ABRIR_CORCH, nil)
	t := p.lookahead.Kind
	if !(t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN ||
		t == token.CERRAR_CORCH) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_SENT, nil) // 10")
		return
	}
	p.FuncBody()
	if p.lookahead.Kind != token.CERRAR_CORCH {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_CERRAR_CORCH, nil) // 10")
		return
	}
	p.match(token.CERRAR_CORCH, nil)
	if p.lookahead.Kind != token.WHILE {
		errors.NewError(errors.SINTACTICAL, errors.S_MISSING_WHILE, nil) // 10")
		return
	}
	p.match(token.WHILE, nil)
	if p.lookahead.Kind != token.ABRIR_PAR {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_ABRIR_PAR, nil) // 10")
		return
	}
	p.match(token.ABRIR_PAR, nil)
	t = p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 10")
		return
	}
	p.Expr()
	if p.lookahead.Kind != token.CERRAR_PAR {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 10")
		return
	}
	p.match(token.CERRAR_PAR, nil)
	if p.lookahead.Kind != token.PUNTOYCOMA {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_SEMICOLON, nil) // 10")
		return
	}
	p.match(token.PUNTOYCOMA, nil)
}

func (p *ParserExec) Expr() {
	p.rule(11)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 11")
		return
	}

	p.ExpRel()
	t = p.lookahead.Kind
	if !(t == token.ARITM || t == token.LOGICO || t == token.CERRAR_PAR ||
		t == token.PUNTOYCOMA || t == token.COMA) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 11b")
		return
	}
	p.Expr2()
}

func (p *ParserExec) Expr2() {
	t := p.lookahead.Kind
	if t == token.LOGICO && p.lookahead.Attr == token.LOG_AND {
		p.rule(12)
		p.match(token.LOGICO, token.LOG_AND)
	} else if t == token.ARITM || t == token.LOGICO || t == token.CERRAR_PAR ||
		t == token.PUNTOYCOMA || t == token.COMA {
		p.rule(13) //lambda
	} else {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 12/13")
		return
	}
}

func (p *ParserExec) ExpRel() {
	p.rule(14)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 14")
		return
	}

	p.AritExp()
	t = p.lookahead.Kind
	if !(t == token.RELAC || t == token.LOGICO || t == token.ARITM || t == token.COMA ||
		t == token.PUNTOYCOMA || t == token.CERRAR_PAR) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 14")
		return
	}
	p.ExpRel2()
}
func (p *ParserExec) ExpRel2() {
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
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 15/16")
			return
		}
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 15/16")
			return
		}
		p.AritExp()
		t = p.lookahead.Kind
		if !(t == token.RELAC || t == token.LOGICO || t == token.ARITM || t == token.COMA ||
			t == token.PUNTOYCOMA || t == token.CERRAR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 15/16b")
			return
		}
		p.ExpRel2()
	case token.ARITM, token.COMA, token.PUNTOYCOMA, token.LOGICO, token.CERRAR_PAR:
		p.rule(17)
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 17")
	}
}

func (p *ParserExec) AritExp() {
	p.rule(18)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 18/a")
		return
	}
	p.Term()
	t = p.lookahead.Kind
	if !(t == token.LOGICO || t == token.ARITM || t == token.COMA ||
		t == token.PUNTOYCOMA || t == token.CERRAR_PAR || t == token.RELAC) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 18b")
		return
	}
	p.AritExp2()
}

func (p *ParserExec) AritExp2() {
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
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 20")
			return
		}
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 19")
			return
		}
		p.Term()
		t = p.lookahead.Kind
		if !(t == token.LOGICO || t == token.ARITM || t == token.COMA ||
			t == token.PUNTOYCOMA || t == token.CERRAR_PAR || t == token.RELAC) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 19/20")
			return
		}
		p.AritExp2()
	case token.RELAC, token.LOGICO, token.CERRAR_PAR, token.PUNTOYCOMA, token.COMA:
		p.rule(21) //lambda
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil)
	}
}

func (p *ParserExec) Term() {
	switch p.lookahead.Kind {
	case token.LOGICO:
		if p.lookahead.Attr != token.LOG_NEG {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 22")
			return
		}
		p.rule(22)
		p.match(token.LOGICO, token.LOG_NEG)
		t := p.lookahead.Kind
		if !(t == token.TRUE || t == token.FALSE || t == token.INT_LITERAL ||
			t == token.REAL_LITERAL || t == token.STRING_LITERAL ||
			t == token.ABRIR_PAR || t == token.ID) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 22")
			return
		}
		p.Term3()
	case token.ARITM:
		switch p.lookahead.Attr {
		case token.ARIT_PLUS:
			p.rule(23)
			p.match(token.ARITM, token.ARIT_PLUS)
		case token.ARIT_MINUS:
			p.rule(24)
			p.match(token.ARITM, token.ARIT_MINUS)
		default:
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 23/24")
			return
		}
		switch p.lookahead.Kind {
		case token.INT_LITERAL, token.REAL_LITERAL, token.ID, token.STRING_LITERAL,
			token.ABRIR_PAR:
			p.Term2()
		default:
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 23/24")
		}
	case token.INT_LITERAL, token.REAL_LITERAL, token.ID, token.STRING_LITERAL,
		token.ABRIR_PAR:
		p.rule(25)
		p.Term2()
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 23/24")
	}
}

func (p *ParserExec) Term3() {
	switch p.lookahead.Kind {
	case token.TRUE:
		p.rule(26)
		p.match(token.TRUE, nil)
	case token.FALSE:
		p.rule(27)
		p.match(token.FALSE, nil)
	case token.INT_LITERAL, token.REAL_LITERAL, token.ID, token.STRING_LITERAL,
		token.ABRIR_PAR:
		p.rule(28)
		p.Term2()
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 28")
		return
	}
}

func (p *ParserExec) Term2() {
	switch p.lookahead.Kind {
	case token.INT_LITERAL:
		p.rule(29)
		p.match(token.INT_LITERAL, nil)
	case token.REAL_LITERAL:
		p.rule(30)
		p.match(token.REAL_LITERAL, nil)
	case token.ID:
		p.rule(31)
		p.match(token.ID, nil)
		t := p.lookahead.Kind
		if !(t == token.ABRIR_PAR || t == token.CERRAR_PAR || t == token.ARITM ||
			t == token.RELAC || t == token.LOGICO || t == token.COMA ||
			t == token.PUNTOYCOMA) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 31")
			return
		}
		p.FactorId()
	case token.STRING_LITERAL:
		p.rule(32)
		p.match(token.STRING_LITERAL, nil)
	case token.ABRIR_PAR:
		p.rule(33)
		p.match(token.ABRIR_PAR, nil)
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 33")
			return
		}
		p.Expr()
		if p.lookahead.Kind != token.CERRAR_PAR {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 33")
			return
		}
		p.match(token.CERRAR_PAR, nil)
	}
}

func (p *ParserExec) FactorId() {
	switch p.lookahead.Kind {
	case token.ABRIR_PAR:
		p.rule(34)
		p.match(token.ABRIR_PAR, nil)
		switch p.lookahead.Kind {
		case token.ARITM, token.LOGICO, token.CERRAR_PAR, token.ID, token.INT_LITERAL,
			token.REAL_LITERAL, token.STRING_LITERAL:
			if p.lookahead.Kind == token.LOGICO {
				if p.lookahead.Attr != token.LOG_NEG {
					errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 34a")
					return
				}
			}
			p.ParamList()
		default:
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 34b")
			return
		}
		if p.lookahead.Kind != token.CERRAR_PAR {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 34c")
			return
		}
		p.match(token.CERRAR_PAR, nil)

	case token.LOGICO, token.ARITM, token.RELAC, token.COMA, token.PUNTOYCOMA,
		token.CERRAR_PAR:
		p.rule(35)
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 35")
	}
}

func (p *ParserExec) DecFunc() {
	p.rule(36)
	if p.lookahead.Kind != token.FUNCTION {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.match(token.FUNCTION, nil)
	t := p.lookahead.Kind
	if !(t == token.STRING || t == token.VOID || t == token.INT ||
		t == token.FLOAT || t == token.BOOLEAN) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.TipoFunc()
	if p.lookahead.Kind != token.ID {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.match(token.ID, nil)
	if p.lookahead.Kind != token.ABRIR_PAR {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.match(token.ABRIR_PAR, nil)
	t = p.lookahead.Kind
	if !(t == token.STRING || t == token.VOID || t == token.INT ||
		t == token.FLOAT || t == token.BOOLEAN) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.FuncParams()
	if p.lookahead.Kind != token.CERRAR_PAR {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.match(token.CERRAR_PAR, nil)
	if p.lookahead.Kind != token.ABRIR_CORCH {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.match(token.ABRIR_CORCH, nil)
	t = p.lookahead.Kind
	if !(t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN ||
		t == token.CERRAR_CORCH) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.FuncBody()
	if p.lookahead.Kind != token.CERRAR_CORCH {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 36")
		return
	}
	p.match(token.CERRAR_CORCH, nil)
}

func (p *ParserExec) TipoFunc() {
	switch p.lookahead.Kind {
	case token.VOID:
		p.rule(38)
		p.match(token.VOID, nil)
	case token.STRING, token.INT, token.FLOAT, token.BOOLEAN:
		p.rule(37)
		p.Tipo()
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 37/38")
		return
	}
}

func (p *ParserExec) FuncParams() {
	switch p.lookahead.Kind {
	case token.INT, token.FLOAT, token.BOOLEAN, token.STRING: //FIRST(Tipo)
		p.rule(39)
		p.Tipo()
		if p.lookahead.Kind != token.ID {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 39")
			return
		}
		p.match(token.ID, nil)
		if p.lookahead.Kind == token.COMA || p.lookahead.Kind == token.CERRAR_PAR {
			//first FuncParams2
			p.FuncParams2()
		} else {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 39")
			return
		}
	case token.VOID:
		p.match(token.VOID, nil)
		p.rule(40)
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 40")
	}
}

func (p *ParserExec) FuncParams2() {
	switch p.lookahead.Kind {
	case token.COMA:
		p.match(token.COMA, nil)
		p.rule(41)
		switch p.lookahead.Kind {
		case token.INT, token.FLOAT, token.BOOLEAN, token.STRING, token.ID: //FIRST(Tipo)&id
			p.Tipo()
			if p.lookahead.Kind != token.ID {
				errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 41a")
				return
			}
			p.match(token.ID, nil)
			if !(p.lookahead.Kind == token.COMA || p.lookahead.Kind == token.CERRAR_PAR) {
				errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 41b")
				return
			}
			p.FuncParams2()
		default:
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 41")
			return
		}
	case token.CERRAR_PAR:
		p.rule(42)
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 42")
		return
	}
}

func (p *ParserExec) Tipo() {
	switch p.lookahead.Kind {
	case token.INT:
		p.match(token.INT, nil)
		p.rule(43)
	case token.FLOAT:
		p.match(token.FLOAT, nil)
		p.rule(44)
	case token.BOOLEAN:
		p.match(token.BOOLEAN, nil)
		p.rule(45)
	case token.STRING:
		p.match(token.STRING, nil)
		p.rule(46)
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 43-46")
	}
}

func (p *ParserExec) FuncBody() {
	switch p.lookahead.Kind {
	case token.IF, token.LET, token.DO, token.ID, token.READ, token.WRITE,
		token.RETURN: //first Decl
		p.rule(47)
		p.Decl()
		t := p.lookahead.Kind
		if !(t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
			t == token.WRITE || t == token.READ || t == token.RETURN ||
			t == token.CERRAR_CORCH) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 47")
			return
		}
		p.FuncBody()

	case token.CERRAR_CORCH:
		p.rule(48)

	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 48")
		return
	}
}

func (p *ParserExec) ParamList() {
	t := p.lookahead.Kind
	if t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR {
		p.rule(49)
		p.Expr()
		t = p.lookahead.Kind
		if !(t == token.COMA || t == token.CERRAR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 49")
			return
		}
		p.ParamList2()

	} else if p.lookahead.Kind == token.CERRAR_PAR {
		p.rule(50)
	} else {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 50")
		return
	}
}

func (p *ParserExec) ParamList2() {
	switch p.lookahead.Kind {
	case token.COMA:
		p.match(token.COMA, nil)
		p.rule(51)
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 51")
			return
		}
		p.Expr()
		t = p.lookahead.Kind
		if !(t == token.COMA || t == token.CERRAR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 51")
			return
		}
		p.ParamList2()
	case token.CERRAR_PAR: //FOLLOW={)
		p.rule(52)
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 52")
	}
}

func (p *ParserExec) Sent() {
	switch p.lookahead.Kind {
	case token.ID:
		p.match(token.ID, nil)
		p.rule(53)
		switch p.lookahead.Kind {
		case token.ASIG, token.ABRIR_PAR:
			break
		default:
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 53")
			return
		}
		p.Sent2()

	case token.WRITE:
		p.match(token.WRITE, nil)
		p.rule(54)
		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 54")
			return
		}
		p.Expr()
		p.match(token.PUNTOYCOMA, nil)

	case token.READ:
		p.rule(55)
		if !p.match(token.READ, nil) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 55")
			return
		}
		if !p.match(token.ID, nil) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 55")
			return
		}
		if !p.match(token.PUNTOYCOMA, nil) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 55")
			return
		}

	case token.RETURN:
		p.rule(56)
		if !p.match(token.RETURN, nil) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 56")
			return
		}
		p.ReturnExp()
		if !p.match(token.PUNTOYCOMA, nil) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 55")
			return
		}

	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 53-56")
	}

}
func (p *ParserExec) Sent2() {
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
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 57/58")
			return
		}

		t := p.lookahead.Kind
		if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
			t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
			t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 57/58")
			return
		}
		p.Expr()
		if p.lookahead.Kind != token.PUNTOYCOMA {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 57")
			return
		}
		p.match(token.PUNTOYCOMA, nil)

	case token.ABRIR_PAR:
		p.match(token.ABRIR_PAR, nil)
		p.rule(59)
		switch p.lookahead.Kind {
		case token.ARITM, token.LOGICO, token.CERRAR_PAR, token.ID, token.INT_LITERAL,
			token.REAL_LITERAL, token.STRING_LITERAL:
			if p.lookahead.Kind == token.LOGICO {
				if p.lookahead.Attr != token.LOG_NEG {
					errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 59")
					return
				}
			}
		default:
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 59")
			return
		}

		p.ParamList()
		if p.lookahead.Kind != token.CERRAR_PAR {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 59")
			return
		}
		p.match(token.CERRAR_PAR, nil)

		if p.lookahead.Kind != token.PUNTOYCOMA {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 59")
			return
		}
		p.match(token.PUNTOYCOMA, nil)
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 57-59")
	}
}
func (p *ParserExec) ReturnExp() {
	switch p.lookahead.Kind {
	case token.LOGICO:
		if p.lookahead.Attr != token.LOG_NEG {
			errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 60/61")
			return
		}
	case token.ARITM, token.INT_LITERAL, token.REAL_LITERAL, token.STRING_LITERAL,
		token.ID, token.ABRIR_PAR:
		break
	case token.PUNTOYCOMA:
		p.rule(61)
		return
	default:
		errors.NewError(errors.SINTACTICAL, errors.S_INVALID_EXP, nil) // 61")
		return
	}
	p.rule(60)
	t := p.lookahead.Kind
	if !(t == token.ARITM || (t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG) ||
		t == token.INT_LITERAL || t == token.REAL_LITERAL || t == token.ID ||
		t == token.STRING_LITERAL || t == token.ABRIR_PAR) {
		errors.NewError(errors.SINTACTICAL, errors.S_EXPECTED_EXP, nil) // 60")
		return
	}
	p.Expr()

}
