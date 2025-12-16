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

func (p *Parser) hasTokens() bool {
	nextTokenIndex := p.currentTokenIndex + 1
	return nextTokenIndex < len(p.tokens)
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

// AST := FUNCTION_DEF*
func (p *Parser) Parse() (ast.Ast, error) {

	functionDefs := make(map[string]ast.FunctionDef, 0)

	for p.hasTokens() {
		functionDef, err := p.parseFunctionDef()
		if err != nil {
			return ast.Ast{}, err
		}

		if _, ok := functionDefs[functionDef.Name]; ok {
			return ast.Ast{}, fmt.Errorf("Function %+v redefined", functionDef.Name)
		}

		functionDefs[functionDef.Name] = functionDef
	}

	main, ok := functionDefs["main"]
	if !ok {
		return ast.Ast{}, fmt.Errorf("main function missing")
	}

	if len(main.Params) > 0 {
		return ast.Ast{}, fmt.Errorf("main function does not accept any parameters")
	}

	ast := ast.Ast{
		FunctionDefs: functionDefs,
	}

	return ast, nil
}

// FUNC_DEF :=  'def' IDENT '(' NON_EMPTY_FUNCTION_PARAMETER_LIST? ')' statement_block
func (p *Parser) parseFunctionDef() (ast.FunctionDef, error) {
	p.expect(tokenizer.Token_Def)

	name := p.expect(tokenizer.Token_Identifier)
	typedName := name.(tokenizer.TokenWithData)

	p.expect(tokenizer.Token_LParen)

	parameters := make([]ast.Var, 0)

	if !p.nextTokenIs(tokenizer.Token_RParen) {
		ps, err := p.parseFunctionParameterList()
		if err != nil {
			return ast.FunctionDef{}, err
		}
		parameters = ps
	}

	p.expect(tokenizer.Token_RParen)

	statements, err := p.parseStatementBlock()
	if err != nil {
		return ast.FunctionDef{}, err
	}

	return ast.FunctionDef{
		Name:     typedName.Value(),
		Params:   parameters,
		StmtList: statements,
	}, nil
}

// NON_EMPTY_FUNCTION_PARAMETER_LIST = Ident {',' Ident}*
func (p *Parser) parseFunctionParameterList() ([]ast.Var, error) {
	parameters := make([]ast.Var, 0)

	paramToken := p.expect(tokenizer.Token_Identifier)
	typedParamToken := paramToken.(tokenizer.TokenWithData)

	param := ast.Var{
		Var: typedParamToken.Value(),
	}

	parameters = append(parameters, param)

	for p.nextTokenIs(tokenizer.Token_Comma) {
		p.expect(tokenizer.Token_Comma)
		paramToken := p.expect(tokenizer.Token_Identifier)
		typedParamToken := paramToken.(tokenizer.TokenWithData)

		param := ast.Var{
			Var: typedParamToken.Value(),
		}

		parameters = append(parameters, param)

	}

	return parameters, nil
}

// STATEMENT_BLOCK = '{' {PRINT_STATEMENT | ASSIGN_STATEMENT | RETURN {EXPR}?}* '}'
func (p *Parser) parseStatementBlock() ([]ast.Statement, error) {
	statements := make([]ast.Statement, 0)
	p.expect(tokenizer.Token_LBrace)

	for !p.nextTokenIs(tokenizer.Token_RBrace) {

		if p.nextTokenIs(tokenizer.Token_Print) {
			stmt, err := p.statementPrint()
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
			continue
		}

		if p.nextTokenIs(tokenizer.Token_Return) {
			stmt, err := p.statementReturn()
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
			continue
		}

		// otherwise assign statement
		stmt, err := p.statementAssign()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	p.expect(tokenizer.Token_RBrace)
	return statements, nil
}

// STATEMENT_RETURN := 'return' {EXPR}?
func (p *Parser) statementReturn() (ast.Statement, error) {
	p.expect(tokenizer.Token_Return)

	if p.nextTokenIs(tokenizer.Token_Semi) {
		p.expect(tokenizer.Token_Semi)

		return ast.ReturnStmt{}, nil
	}

	// otherwise we have an expression
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.expect(tokenizer.Token_Semi)

	return ast.ReturnWithExprStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) statementAssign() (ast.Statement, error) {
	tok := p.expect(tokenizer.Token_Identifier)

	tokIdent, _ := tok.(tokenizer.TokenWithData)
	identValue := tokIdent.Value()

	p.expect(tokenizer.Token_Assign)

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.expect(tokenizer.Token_Semi)

	return ast.AssignStmt{
		Variable: identValue,
		Expr:     expr,
	}, nil
}

func (p *Parser) statementPrint() (ast.Statement, error) {
	p.expect(tokenizer.Token_Print)

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.expect(tokenizer.Token_Semi)

	return ast.PrintStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
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

// EXPR_MUL := EXPR_UNARY { {* | /} EXPR_UNARY}*
func (p *Parser) expressionMultiplicative() (ast.Expr, error) {
	lhs, err := p.expressionUnary()
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

		rhs, err := p.expressionUnary()
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

// EXPR_UNARY: {`-`}? EXPR_ATOM
func (p *Parser) expressionUnary() (ast.Expr, error) {
	if p.nextTokenIs(tokenizer.Token_Minus) {
		// TODO BOOLEAN Negation.
		_ = p.expect(tokenizer.Token_Minus)

		expr, err := p.expressionAtom()
		if err != nil {
			return nil, err
		}

		return ast.UnaryExpr{
			Op:   ast.UnaryOp_Neg,
			Expr: expr,
		}, nil
	}

	return p.expressionAtom()
}

// EXPR_ATOM := int | identifier |
// EXPR_ATOM := '(' EXPR ')'
// EXPR_ATOM := identifier '(' {EXPR ','} * ')'
// EXPR_ATOM = '&' identifier
func (p *Parser) expressionAtom() (ast.Expr, error) {
	if p.nextTokenIs(tokenizer.Token_FuncAddr) {
		p.expect(tokenizer.Token_FuncAddr)

		tok, ok := p.expect(tokenizer.Token_Identifier).(tokenizer.TokenWithData)
		if !ok {
			panic("bad")
		}

		return ast.AddressOfFunction{
			Name: tok.Value(),
		}, nil
	}

	if p.nextTokenIs(tokenizer.Token_LParen) {
		_ = p.expect(tokenizer.Token_LParen)

		expr, err := p.expressionAdditive()
		if err != nil {
			return nil, err
		}

		_ = p.expect(tokenizer.Token_RParen)

		return expr, nil
	}

	if p.nextTokenIs(tokenizer.Token_Identifier) {
		token, ok := p.expect(tokenizer.Token_Identifier).(tokenizer.TokenWithData)
		if !ok {
			panic("bad")
		}

		// function call
		if p.nextTokenIs(tokenizer.Token_LParen) {
			expressionList := []ast.Expr{}
			p.expect(tokenizer.Token_LParen)
			if !p.nextTokenIs(tokenizer.Token_RParen) {
				e, err := p.expressionList()
				if err != nil {
					return nil, err
				}
				expressionList = e
			}
			p.expect(tokenizer.Token_RParen)

			return ast.FunctionCall{
				Name: token.Value(),
				Args: expressionList,
			}, nil
		}

		return ast.Var{
			Var: token.Value(),
		}, nil
	}

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

// NON-empty expression list
func (p *Parser) expressionList() ([]ast.Expr, error) {
	expressions := make([]ast.Expr, 0)
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	expressions = append(expressions, expr)

	for p.nextTokenIs(tokenizer.Token_Comma) {
		p.expect(tokenizer.Token_Comma)
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}
