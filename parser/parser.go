package parser

import (
	"compiler/ast"
	"compiler/tokenizer"
	"fmt"
	"strconv"
)

type Parser struct {
	tokens            []tokenizer.Token
	currentTokenIndex int
}

func NewParser(tokens []tokenizer.Token) (*Parser, error) {
	if tokens == nil || len(tokens) == 0 {
		return nil, fmt.Errorf("Nothing to parse")
	}

	return &Parser{
		tokens:            tokens,
		currentTokenIndex: -1,
	}, nil
}

func (p *Parser) nextTokenIs(expectedTokenKinds ...tokenizer.TokenKind) bool {
	nextTokenIndex := p.currentTokenIndex + 1

	if nextTokenIndex >= len(p.tokens) {
		return false
	}

	nextToken := p.tokens[nextTokenIndex]
	for _, tokenKind := range expectedTokenKinds {

		if nextToken.Kind() == tokenKind {
			return true
		}

	}

	return false
}

func (p *Parser) expect(expectedTokenKinds ...tokenizer.TokenKind) tokenizer.Token {
	nextTokenIndex := p.currentTokenIndex + 1
	if nextTokenIndex >= len(p.tokens) {
		panic("Reached end of token list")
	}

	nextToken := p.tokens[nextTokenIndex]
	for _, tokenKind := range expectedTokenKinds {

		if nextToken.Kind() == tokenKind {
			p.currentTokenIndex = nextTokenIndex
			return nextToken
		}

	}

	panic(fmt.Sprintf("Error: expected token %+v, actual %+v", expectedTokenKinds, nextToken.Kind()))
}

func (p *Parser) Parse() (ast.Expr, error) {
	return p.expressionAdditive()
}

// EXPR_ADD := EXPR_MUL { {+ | -} EXPR_MUL}*
func (p *Parser) expressionAdditive() (ast.Expr, error) {
	lhs, err := p.expressionMultiplicative()
	if err != nil {
		return nil, err
	}

	addOps := []tokenizer.TokenKind{tokenizer.Token_Plus, tokenizer.Token_Minus}

	for p.nextTokenIs(addOps...) {
		opToken := p.expect(addOps...)
		var op ast.BinOp
		switch opToken.Kind() {
		case tokenizer.Token_Plus:
			op = ast.BinOp_Add
		case tokenizer.Token_Minus:
			op = ast.BinOp_Sub
		}

		rhs, err := p.expressionMultiplicative()
		if err != nil {
			return nil, err
		}

		lhs = ast.BinaryExpr{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		}
	}

	return lhs, nil
}

// EXPR_MUL := EXPR_ATOM { {* | /} EXPR_ATOM}*
func (p *Parser) expressionMultiplicative() (ast.Expr, error) {
	lhs, err := p.expressionAtom()
	if err != nil {
		return nil, err
	}

	ops := []tokenizer.TokenKind{tokenizer.Token_Mul, tokenizer.Token_Div}

	for p.nextTokenIs(ops...) {
		opToken := p.expect(ops...)
		var op ast.BinOp
		switch opToken.Kind() {
		case tokenizer.Token_Mul:
			op = ast.BinOp_Mul
		case tokenizer.Token_Div:
			op = ast.BinOp_Div
		}

		rhs, err := p.expressionAtom()
		if err != nil {
			return nil, err
		}

		lhs = ast.BinaryExpr{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		}
	}

	return lhs, nil

}

// EXPR_ATOM := int
func (p *Parser) expressionAtom() (ast.Expr, error) {
	if p.nextTokenIs(tokenizer.Token_Int) {
		token, ok := p.expect(tokenizer.Token_Int).(tokenizer.TokenWithData)
		if !ok {
			panic("bad")
		}

		parsedInt, err := strconv.Atoi(token.Value())
		if err != nil {
			return nil, err
		}

		return ast.Int{
			Integer: parsedInt,
		}, nil
	}

	return nil, fmt.Errorf("Unable to parse EXPR_ATOM")
}
