package parser

import "compiler-pdl/src/lexer"
import "compiler-pdl/src/token"

type Parser struct {
	lexer     lexer.Lexer
	lookahead token.Token
}

func (p *Parser) askToken() bool {
	tk, ok := p.lexer.Lexical()
	if ok {
		p.lookahead = tk

	}
	return ok
}

// Grammar

func (p *Parser) S() {
	ok := p.askToken()
	if ok {
		switch p.lookahead.Kind {
		case token.FUNCTION:
			{
				p.Decfunc()
			}
		case token.LET, token.ID, token.IF, token.DO, token.READ, token.WRITE:
			{
				p.Sent()
			}
		default: //ERROR
		}

		p.S()
	}
}

func (p *Parser) Decfunc() {
	ok := p.askToken()
	if !ok {
		//ERROR
		return
	}
	p.Tipo()
	if p.lookahead.Kind != token.ID {
		//ERROR
	}

	ok = p.askToken()
	if !ok || p.lookahead.Kind != token.ABRIR_PAR {
		//ERROR
		return
	}

	ok = p.askToken()
	if !ok {
		//ERROR
		return
	}

	p.Exprlist()
	if p.lookahead.Kind != token.CERRAR_PAR {
		//ERROR
	}

	ok = p.askToken()
	if !ok || p.lookahead.Kind != token.ABRIR_CORCH {
		//ERROR
		return
	}
	p.Funcbody()
	ok = p.askToken()
	if !ok || p.lookahead.Kind != token.CERRAR_CORCH {
		//ERROR
		return
	}
}

func (p *Parser) Tipo() {
	switch p.lookahead.Kind {
	case token.STRING, token.INT, token.FLOAT, token.VOID, token.BOOLEAN:
		{
			p.askToken()
		}
	default: // Not Tipo present
	}

}
func (p *Parser) Exprlist() {

}

func (p *Parser) Funcbody(){

}

func (p *Parser) Sent() {

	switch p.lookahead.Kind {
	case token.LET:
		{
		}
	case token.ID:
		{
		}
	case token.IF:
		{
		}
	case token.DO:
		{
		}
	case token.READ:
		{
		}
	case token.WRITE:
		{

		}
	default: //ERROR
	}

}
