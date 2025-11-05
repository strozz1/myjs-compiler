package parser

import "compiler-pdl/src/lexer"
import "compiler-pdl/src/token"
import "compiler-pdl/src/errors"

type Parser struct {
	lexer     lexer.Lexer
	lookahead token.Token
}

func (p *Parser) getToken() bool {
	tk, ok := p.lexer.Lexical()
	if ok {
		p.lookahead = tk

	}
	return ok
}
func (p *Parser) match(tk token.TokenKind, attr any) bool {
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

// Grammar
func (p *Parser) P() {
	switch p.lookahead.Kind {
	case token.FUNCTION: //FIRST(DecFunc)
		p.DecFunc()
	case token.LET, token.ID, token.IF, token.DO, token.READ, token.WRITE, token.RETURN:
		p.Decl()
	default:
		//TODO no mas tokens o error
		return
	}
	t := p.lookahead.Kind
	if !(t == token.FUNCTION || t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN) { //FIRST(P)
		//error
		return
	}
	p.P()
}

func (p *Parser) Decl() {
	switch p.lookahead.Kind {
	case token.IF:
		p.match(token.IF, nil)
		if p.lookahead.Kind != token.ABRIR_PAR {
			//error
			return
		}
		p.match(token.ABRIR_PAR, nil)
		t := p.lookahead.Kind
		if !(t == token.ARITM ||
			(t == token.LOGICO && p.lookahead.Attr == token.LOG_NEG)) {
			//if not FIRST(Expr) error
			return
		}
		p.Expr()
		if p.lookahead.Kind != token.CERRAR_PAR {
			//error
			return
		}
		p.match(token.CERRAR_PAR, nil)
		t = p.lookahead.Kind
		if !(t == token.ID || t == token.WRITE || t == token.READ || t == token.RETURN) {
			//if not FIRST(Sent) error
			return
		}
		p.Sent()
	case token.LET:
		p.match(token.LET, nil)
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.TipoDecl()
		if p.lookahead.Kind != token.ID {
			//error
			return
		}
		p.match(token.ID, nil)
		if p.lookahead.Kind != token.PUNTOYCOMA {
			//error
			return
		}
		p.match(token.PUNTOYCOMA, nil)
	case token.DO:
		p.match(token.DO, nil)
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.WhileBody()
	default:
		t := p.lookahead.Kind
		if t == token.ID || t == token.WRITE || t == token.READ || t == token.RETURN {
			//if FIRST(Sent)
			p.Sent()
		} else {
			//error
			return
		}
	}

}

func (p *Parser) TipoDecl() {
	if p.lookahead.Kind == 0 { //first
		p.Tipo()
	} else { //lambda next
	}
	//error
}

func (p *Parser) WhileBody() {
	if p.lookahead.Kind != token.ABRIR_CORCH {
		//error
		return
	}
	p.match(token.ABRIR_CORCH, nil)
	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.FuncBody()
	if p.lookahead.Kind != token.CERRAR_CORCH {
		//error
		return
	}
	p.match(token.CERRAR_CORCH, nil)
	if p.lookahead.Kind != token.WHILE {
		//error
		return
	}
	p.match(token.WHILE, nil)
	if p.lookahead.Kind != token.ABRIR_PAR {
		//error
		return
	}
	p.match(token.ABRIR_PAR, nil)
	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.Expr()
	if p.lookahead.Kind != token.CERRAR_PAR {
		//error
		return
	}
	p.match(token.CERRAR_PAR, nil)
	if p.lookahead.Kind != token.PUNTOYCOMA {
		//error
		return
	}
	p.match(token.PUNTOYCOMA, nil)
}

func (p *Parser) Expr() {
	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.ExpRel()
	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.Expr2()
}

func (p *Parser) Expr2() {
	if p.lookahead.Kind == token.LOGICO && p.lookahead.Attr == token.LOG_AND {
		p.match(token.LOGICO, token.LOG_AND)
	} else { //lambda or error
	}

}

func (p *Parser) ExpRel() {
	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.AritExp()
	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.ExpRel2()
}
func (p *Parser) ExpRel2() {
	switch p.lookahead.Kind {
	case token.RELAC:
		switch p.lookahead.Attr {
		case token.REL_EQ:
			p.match(token.RELAC, token.REL_EQ)
		case token.REL_NOTEQ:
			p.match(token.RELAC, token.REL_NOTEQ)
		default: //error
			return
		}
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.AritExp()
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.ExpRel2()
	default: //lambda check

	}
}

func (p *Parser) AritExp() {
	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.Term()

	if p.lookahead.Kind != 0 { //first
		//error
		return
	}
	p.AritExp2()
}

func (p *Parser) AritExp2() {
	switch p.lookahead.Kind {
	case token.ARITM:
		switch p.lookahead.Attr {
		case token.ARIT_MINUS:
			p.match(token.ARITM, token.ARIT_MINUS)
		case token.ARIT_PLUS:
			p.match(token.ARITM, token.ARIT_PLUS)
		default: //error
			return
		}
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.AritExp2()
	default: //follow lambda
		//TODO
	}
}

func (p *Parser) Term() {
	switch p.lookahead.Kind {
	case token.LOGICO:
		if p.lookahead.Attr != token.LOG_NEG {
			//errro
			return
		}
		p.match(token.LOGICO, token.LOG_NEG)
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.Term3()
	case token.ARITM:
		switch p.lookahead.Attr {
		case token.ARIT_PLUS:
			p.match(token.ARITM, token.ARIT_PLUS)
		case token.ARIT_MINUS:
			p.match(token.ARITM, token.ARIT_MINUS)
		default:
			//errro
			return
		}
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.Term2()
	}

}

func (p *Parser) Term3() {
	switch p.lookahead.Kind {
	case token.TRUE:
		p.match(token.TRUE, nil)
	case token.FALSE:
		p.match(token.FALSE, nil)
	default:
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.Term2()
	}

}

func (p *Parser) Term2() {
	switch p.lookahead.Kind {
	case token.INT_LITERAL:
		p.match(token.INT_LITERAL, nil)
	case token.REAL_LITERAL:
		p.match(token.REAL_LITERAL, nil)
	case token.ID:
		p.match(token.ID, nil)
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.FactorId()
	case token.STRING_LITERAL:
		p.match(token.STRING_LITERAL, nil)
	case token.ABRIR_PAR:
		p.match(token.ABRIR_PAR, nil)
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.Expr()
		if p.lookahead.Kind != token.CERRAR_PAR {
			//error
			return
		}
		p.match(token.CERRAR_PAR, nil)
	}
}

func (p *Parser) FactorId() {
	if p.lookahead.Kind == token.ABRIR_PAR {
		p.match(token.ABRIR_PAR, nil)
		if p.lookahead.Kind != 0 { //first
			//error
			return
		}
		p.ParamList()
		if p.lookahead.Kind != token.CERRAR_PAR {
			//errr
			return
		}
		p.match(token.CERRAR_PAR, nil)
	} else {
		//lambda todo

	}
}
func (p *Parser) Tipo() {
	t := p.lookahead.Kind
	if t == token.STRING || t == token.INT || t == token.FLOAT || t == token.BOOLEAN {
		p.match(t, nil)
		return
	}
	//error
}
func (p *Parser) DecFunc() {
	if p.lookahead.Kind != token.FUNCTION {
		//error
		return
	}
	p.match(token.FUNCTION, nil)
	t := p.lookahead.Kind
	if !(t == token.STRING || t == token.VOID || t == token.INT || t == token.FLOAT || t == token.BOOLEAN) {
		//if not FIRST(TipoFunc) error
		return
	}
	p.TipoFunc()
	if p.lookahead.Kind != token.ID {
		//error
		return
	}
	p.match(token.ID, nil)
	if p.lookahead.Kind != token.ABRIR_PAR {
		//error
		return
	}
	p.match(token.ABRIR_PAR, nil)
	t = p.lookahead.Kind
	if !(t == token.STRING || t == token.VOID || t == token.INT || t == token.FLOAT || t == token.BOOLEAN) {
		//if not FIRST(FuncParams) error
		return
	}
	p.FuncParams()
	if p.lookahead.Kind != token.CERRAR_PAR {
		//error
		return
	}
	p.match(token.CERRAR_PAR, nil)
	if p.lookahead.Kind != token.ABRIR_CORCH {
		//error
		return
	}
	p.match(token.ABRIR_CORCH, nil)
	t = p.lookahead.Kind
	if !(t == token.IF || t == token.LET || t == token.DO || t == token.ID ||
		t == token.WRITE || t == token.READ || t == token.RETURN || t == token.CERRAR_CORCH) {
		//FIRST(FuncBody) & FIRST(})
		//erro
		return
	}
	p.FuncBody()
	if p.lookahead.Kind != token.CERRAR_CORCH {
		//error
		return
	}
	p.match(token.CERRAR_CORCH, nil)
}

func (p *Parser) TipoFunc() {
	switch p.lookahead.Kind {
	case token.VOID:
		p.match(token.VOID, nil)
	case token.STRING, token.INT, token.FLOAT, token.BOOLEAN:
		p.Tipo()
	default:
		//error
		return
	}
}

func (p *Parser) FuncParams() {
	switch p.lookahead.Kind {
	case token.INT, token.FLOAT, token.BOOLEAN, token.STRING: //FIRST(Tipo)
		p.Tipo()
		if p.lookahead.Kind != token.ID {
			//error
			return
		}
		p.match(token.ID, nil)
		if p.lookahead.Kind == token.COMA || p.lookahead.Kind == token.CERRAR_PAR { //first FuncParams2
			p.FuncParams2()
		} else {
			// error
			return
		}
	case token.CERRAR_PAR:
		//follow={)} -> trans lambda
	default:
		//error
		return
	}
}

func (p *Parser) FuncParams2() {
	switch p.lookahead.Kind {
	case token.COMA:
		p.match(token.COMA, nil)
		switch p.lookahead.Kind {
		case token.INT, token.FLOAT, token.BOOLEAN, token.STRING: //FIRST(Tipo)
			p.Tipo()
			if p.lookahead.Kind != token.ID {
				//error
				return
			}
			if !(p.lookahead.Kind == token.COMA || p.lookahead.Kind == token.CERRAR_PAR) {
				//error
				return
			}
			p.FuncParams2()
		default:
			//error
			return
		}
	case token.CERRAR_PAR:
		//follow=)
	default:
		//error
		return
	}
}

func (p *Parser) FuncBody() {
	switch p.lookahead.Kind {
	case token.IF, token.LET, token.DO, token.ID, token.READ, token.WRITE, token.RETURN: //first Decl
		p.Decl()
		if p.lookahead.Kind == 0 { //first Funcbody
			p.FuncBody()
		}
	case token.CERRAR_CORCH:
		//FOLLOW=}

	default:
		//error
		return
	}
}
func (p *Parser) Exprlist() {

}

func (p *Parser) Funcbody() {

}

func (p *Parser) Sent() {

	switch p.lookahead.Kind {
	case token.ID:
		{
			if !p.match(token.ID, nil) {
				//error
				return
			}
			p.Sent2()
		}
	case token.WRITE:
		{
			if !p.match(token.WRITE, nil) {
				//error
				return
			}
			p.Expr()
			if !p.match(token.PUNTOYCOMA, nil) {
				//error
				return
			}
		}
	case token.READ:
		{
			if !p.match(token.READ, nil) {
				//error
				return
			}
			if !p.match(token.ID, nil) {
				//error
				return
			}
			if !p.match(token.PUNTOYCOMA, nil) {
				//error
				return
			}
		}
	case token.RETURN:
		{
			if !p.match(token.RETURN, nil) {
				//error
				return
			}
			p.ReturnExp()
			if !p.match(token.PUNTOYCOMA, nil) {
				//error
				return
			}
		}
	default: //ERROR
	}

}

func (p *Parser) ParamList() {
	if p.lookahead.Attr == 0 { //TODO Exp
		p.Expr()
		p.ParamList2()

	} else if p.lookahead.Kind == token.CERRAR_PAR {
		//FOLLOW(ParamList)={)}
	} else {
		//error
		return
	}

}

func (p *Parser) ParamList2() {

	switch p.lookahead.Kind {
	case token.COMA:
		if !p.match(token.COMA, nil) {
			//error
			return
		}
		p.Expr()
		p.ParamList2()
	case token.CERRAR_PAR: //FOLLOW={)

	default:
		//error
		return
	}
}

func (p *Parser) Sent2() {

	if p.lookahead.Kind == 0 { //TODO Exp
		p.Expr()
		if !p.getToken() || !p.match(token.PUNTOYCOMA, nil) {
			//TODO ERROR
			return
		}

	} else if p.lookahead.Kind == token.ASIG && p.lookahead.Attr == token.ASIG_MULT {
		p.Expr()

		if !p.getToken() || !p.match(token.CERRAR_PAR, nil) {
			//TODO ERROR
			return
		}
	} else if p.lookahead.Kind == token.ABRIR_PAR {
		p.ParamList()
		if !p.getToken() || !p.match(token.CERRAR_PAR, nil) {
			//TODO ERROR
			return
		}
		if !p.getToken() || !p.match(token.PUNTOYCOMA, nil) {
			//TODO ERROR
			return
		}
	} else {
		//ERROR
	}
}
func (p *Parser) ReturnExp() {
	switch p.lookahead.Kind {
	case token.LOGICO:
		if p.lookahead.Attr != token.LOG_NEG {
			//error
			return
		}
	case token.ARITM:
		break
	case token.PUNTOYCOMA:
		//follow={;}
		return
	default:
		//error
		errors.NewError(errors.K_SINTACTICAL, 0, "error")
		return
	}
	p.Expr()

}
