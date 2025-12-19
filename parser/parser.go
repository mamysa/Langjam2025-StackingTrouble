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

func (p *Parser) ensureSymbolNotDefined(functionDefs map[string]ast.FunctionDef, globalDefs map[string]ast.GlobalDef, sym string) error {
	if _, ok := functionDefs[sym]; ok {
		return fmt.Errorf("Symbol %+v redefined", sym)
	}

	if _, ok := globalDefs[sym]; ok {
		return fmt.Errorf("Symbol %+v redefined", sym)
	}

	return nil
}

// AST := FUNCTION_DEF*
func (p *Parser) Parse() (ast.Ast, error) {

	functionDefs := make(map[string]ast.FunctionDef, 0)
	globalDefs := make(map[string]ast.GlobalDef, 0)

	for p.hasTokens() {
		if p.nextTokenIs(tokenizer.Token_Global) {
			globalDef, err := p.parseGlobalDef()
			if err != nil {
				return ast.Ast{}, err
			}

			if err := p.ensureSymbolNotDefined(functionDefs, globalDefs, globalDef.GlobalName); err != nil {
				return ast.Ast{}, err
			}

			globalDefs[globalDef.GlobalName] = globalDef
			continue

		}

		functionDef, err := p.parseFunctionDef()
		if err != nil {
			return ast.Ast{}, err
		}

		if err := p.ensureSymbolNotDefined(functionDefs, globalDefs, functionDef.Name); err != nil {
			return ast.Ast{}, err
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
		GlobalDefs:   globalDefs,
	}

	return ast, nil
}

// global_def := `global` ident = int | float
func (p *Parser) parseGlobalDef() (ast.GlobalDef, error) {
	p.expect(tokenizer.Token_Global)

	globalName := p.expect(tokenizer.Token_Identifier).(tokenizer.TokenWithData)

	p.expect(tokenizer.Token_Assign)

	// quick dirty hack, need to have negative constants.
	sign := 1
	if p.nextTokenIs(tokenizer.Token_Minus) {
		p.expect(tokenizer.Token_Minus)
		sign = -1
	}

	if p.nextTokenIs(tokenizer.Token_Int) {
		constInt := p.expect(tokenizer.Token_Int).(tokenizer.TokenWithData)

		parsedInt, err := strconv.Atoi(constInt.Value())
		if err != nil {
			return ast.GlobalDef{}, err
		}

		return ast.GlobalDef{
			GlobalName: globalName.Value(),
			ValueInt: &ast.Int{
				Integer: parsedInt * sign,
			},
		}, nil

	}

	if p.nextTokenIs(tokenizer.Token_Float) {
		constFloat := p.expect(tokenizer.Token_Float).(tokenizer.TokenWithData)

		parsedFloat, err := strconv.ParseFloat(constFloat.Value(), 64)
		if err != nil {
			return ast.GlobalDef{}, err
		}

		return ast.GlobalDef{
			GlobalName: globalName.Value(),
			ValueFloat: &ast.Float{
				Float: parsedFloat * float64(sign),
			},
		}, nil
	}

	panic("unable to parse global constant definition")
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

		if p.nextTokenIs(tokenizer.Token_If) {
			stmt, err := p.statementIf()
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
			continue
		}

		if p.nextTokenIs(tokenizer.Token_Assert) {
			stmt, err := p.statementAssert()
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
			continue
		}

		if p.nextTokenIs(tokenizer.Token_While) {
			stmt, err := p.statementWhile()
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
			continue
		}

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

// statement_assert = `assert` `(` expression `)`
func (p *Parser) statementAssert() (ast.Statement, error) {
	p.expect(tokenizer.Token_Assert)
	p.expect(tokenizer.Token_LParen)

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.expect(tokenizer.Token_RParen)
	p.expect(tokenizer.Token_Semi)

	return ast.AssertStmt{
		Expr: expr,
	}, nil
}

// STATEMENT_WHILE `while` EXPRESSION STATEMENT_BLOCK
func (p *Parser) statementWhile() (ast.Statement, error) {
	p.expect(tokenizer.Token_While)

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	whileBodyBlock, err := p.parseStatementBlock()
	if err != nil {
		return nil, err
	}

	return ast.WhileStmt{
		Cond: cond,
		Body: whileBodyBlock,
	}, nil
}

// STATEMENT_IF: `if`  EXPRESSION   STATEMENT_BLOCK  { `else` STATEMENT_BLOCK }?
func (p *Parser) statementIf() (ast.Statement, error) {
	p.expect(tokenizer.Token_If)

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	ifStatementBlock, err := p.parseStatementBlock()
	if err != nil {
		return nil, err
	}

	elseStatementBlock := []ast.Statement{}

	if p.nextTokenIs(tokenizer.Token_Else) {
		p.expect(tokenizer.Token_Else)

		b, err := p.parseStatementBlock()
		if err != nil {
			return nil, err
		}
		elseStatementBlock = b
	}

	return ast.IfStmt{
		Cond:       cond,
		IfBranch:   ifStatementBlock,
		ElseBranch: elseStatementBlock,
	}, nil
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
	// assignment expression
	aexpr, err := p.expressionLhsIdentTrailer()
	if err != nil {
		return nil, err
	}

	// void function call
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
		p.expect(tokenizer.Token_Semi)
		return ast.VoidFunctionCall{
			Expr: aexpr,
			Args: expressionList,
		}, err
	}

	//definitely assignment expression
	if p.nextTokenIs(tokenizer.Token_Assign) {
		p.expect(tokenizer.Token_Assign)

		var assignmentExpr ast.AssignmentExpr

		if e, ok := aexpr.(ast.Var); ok {
			assignmentExpr = ast.AssignVar{Var: e.Var}
		}

		if e, ok := aexpr.(ast.SubscriptGet); ok {
			assignmentExpr = ast.SubscriptSet{
				Expr:      e.Expr,
				Subscript: e.Subscript,
			}
		}

		if e, ok := aexpr.(ast.ObjectFieldGet); ok {
			assignmentExpr = ast.ObjectFieldSet{
				Expr:  e.Expr,
				Field: e.Field,
			}
		}

		if assignmentExpr == nil {
			return nil, fmt.Errorf("Unable to match AssignmentExpr for %+v", aexpr)
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		p.expect(tokenizer.Token_Semi)

		return ast.AssignStmt{
			AssignmentExpr: assignmentExpr,
			Expr:           expr,
		}, nil
	}

	panic("unable to parse expression")
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
	return p.expressionTernary()
}

// expression_ternary =  expression_disjunction {'if' expression_disjunction 'else' expression_disjunction}?
func (p *Parser) expressionTernary() (ast.Expr, error) {
	expr, err := p.expressionDisjunction()
	if err != nil {
		return nil, err
	}

	if p.nextTokenIs(tokenizer.Token_If) {
		p.expect(tokenizer.Token_If)

		cond, err := p.expressionDisjunction()
		if err != nil {
			return nil, err
		}

		p.expect(tokenizer.Token_Else)

		elseExpr, err := p.expressionDisjunction()
		if err != nil {
			return nil, err
		}

		return ast.TernaryExpr{
			Cond:     cond,
			ThenExpr: expr,
			ElseExpr: elseExpr,
		}, nil
	}

	return expr, nil
}

// EXPR_DISJUNCTION = expr_conjunction {`or` expr_conjunction} *
func (p *Parser) expressionDisjunction() (ast.Expr, error) {
	lhs, err := p.expressionConjunction()
	if err != nil {
		return nil, err
	}

	for p.nextTokenIs(tokenizer.Token_Or) {
		p.expect(tokenizer.Token_Or)

		rhs, err := p.expressionConjunction()
		if err != nil {
			return nil, err
		}

		lhs = ast.BinaryExpr{
			Op:  ast.BinOp_Or,
			Lhs: lhs,
			Rhs: rhs,
		}
	}

	return lhs, nil
}

// EXPR_CONJUNCTION = expr_not {`or` expr_not} * // for now we go directly to exressionComparison
func (p *Parser) expressionConjunction() (ast.Expr, error) {

	lhs, err := p.expressionNot()
	if err != nil {
		return nil, err
	}

	for p.nextTokenIs(tokenizer.Token_And) {
		p.expect(tokenizer.Token_And)

		rhs, err := p.expressionNot()
		if err != nil {
			return nil, err
		}

		lhs = ast.BinaryExpr{
			Op:  ast.BinOp_And,
			Lhs: lhs,
			Rhs: rhs,
		}
	}

	return lhs, nil
}

func (p *Parser) expressionNot() (ast.Expr, error) {
	if p.nextTokenIs(tokenizer.Token_Not) {
		p.expect(tokenizer.Token_Not)
		expr, err := p.expressionComparison()
		if err != nil {
			return nil, err
		}

		return ast.UnaryExpr{
			Op:   ast.UnaryOp_Not,
			Expr: expr,
		}, nil

	}

	return p.expressionComparison()
}

// EXPR_COMPARISON = expr_add {`<` | `>=` | `!=` | `==` expr_add}?
func (p *Parser) expressionComparison() (ast.Expr, error) {
	lhs, err := p.expressionAdditive()
	if err != nil {
		return nil, err
	}

	expectedTokens := []tokenizer.TokenKind{tokenizer.Token_EqEq, tokenizer.Token_NotEq, tokenizer.Token_Lt, tokenizer.Token_GrEq}
	if p.nextTokenIs(expectedTokens...) {
		token := p.expect(expectedTokens...)
		var op ast.BinOp
		switch token.Kind() {
		case tokenizer.Token_EqEq:
			op = ast.BinOp_EqEq
		case tokenizer.Token_NotEq:
			op = ast.BinOp_NotEq
		case tokenizer.Token_Lt:
			op = ast.BinOp_Lt
		case tokenizer.Token_GrEq:
			op = ast.BinOp_GrEq
		default:
			panic("unknown comparison token")
		}

		rhs, err := p.expressionAdditive()
		if err != nil {
			return nil, err
		}

		return ast.BinaryExpr{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		}, nil
	}

	return lhs, nil

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
// expr_atom = '[' ']'     (new-array-expr)
// expr_atom = '{' '}'     (new-object-expr)
// expr_atom = `len` `(` expr `)`
func (p *Parser) expressionAtom() (ast.Expr, error) {
	if p.nextTokenIs(tokenizer.Token_None) {
		p.expect(tokenizer.Token_None)
		return ast.None{}, nil
	}

	if p.nextTokenIs(tokenizer.Token_True) {
		p.expect(tokenizer.Token_True)
		return ast.Bool{
			Value: true,
		}, nil
	}

	if p.nextTokenIs(tokenizer.Token_False) {
		p.expect(tokenizer.Token_False)
		return ast.Bool{
			Value: false,
		}, nil
	}

	if p.nextTokenIs(tokenizer.Token_Len) {
		p.expect(tokenizer.Token_Len)
		p.expect(tokenizer.Token_LParen)
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.expect(tokenizer.Token_RParen)

		return ast.LenExpr{Expr: expr}, nil
	}

	if p.nextTokenIs(tokenizer.Token_CastFloat) {
		p.expect(tokenizer.Token_CastFloat)
		p.expect(tokenizer.Token_LParen)
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.expect(tokenizer.Token_RParen)

		return ast.CastFloat{Expr: expr}, nil
	}

	if p.nextTokenIs(tokenizer.Token_Global) {
		p.expect(tokenizer.Token_Global)
		p.expect(tokenizer.Token_LParen)
		globalName := p.expect(tokenizer.Token_Identifier).(tokenizer.TokenWithData)
		p.expect(tokenizer.Token_RParen)

		return ast.ReadGlobal{GlobalName: globalName.Value()}, nil
	}

	if p.nextTokenIs(tokenizer.Token_LBracket) {
		p.expect(tokenizer.Token_LBracket)
		p.expect(tokenizer.Token_RBracket)
		return ast.NewListExpr{}, nil
	}

	if p.nextTokenIs(tokenizer.Token_LBrace) {
		p.expect(tokenizer.Token_LBrace)
		p.expect(tokenizer.Token_RBrace)
		return ast.NewObjectExpr{}, nil
	}

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

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_ = p.expect(tokenizer.Token_RParen)

		return expr, nil
	}

	if p.nextTokenIs(tokenizer.Token_Identifier) {
		expr, err := p.expressionRhsIdentTrailer()
		if err != nil {
			return nil, err
		}

		return expr, nil
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

	if p.nextTokenIs(tokenizer.Token_Float) {
		token, ok := p.expect(tokenizer.Token_Float).(tokenizer.TokenWithData)
		if !ok {
			panic("bad")
		}

		parsedFloat, err := strconv.ParseFloat(token.Value(), 64)
		if err != nil {
			return nil, err
		}

		return ast.Float{
			Float: parsedFloat,
		}, nil

	}

	return nil, fmt.Errorf("Unable to parse EXPR_ATOM")
}

// expression_lhs_ident_trailer =   rhs_ident_trailer_list_subscript
// left side of assignment doesnt allow function calls
func (p *Parser) expressionLhsIdentTrailer() (ast.Expr, error) {
	token, ok := p.expect(tokenizer.Token_Identifier).(tokenizer.TokenWithData)
	if !ok {
		panic("bad")
	}

	var expr ast.Expr = ast.Var{
		Var: token.Value(),
	}

	trailerBeginTokens := []tokenizer.TokenKind{tokenizer.Token_LBracket, tokenizer.Token_Dot}
	for p.nextTokenIs(trailerBeginTokens...) {
		if p.nextTokenIs(tokenizer.Token_LBracket) {
			e, err := p.identifierTrailerListSubscript(expr)
			if err != nil {
				return nil, err
			}
			expr = e
		}

		if p.nextTokenIs(tokenizer.Token_Dot) {
			e, err := p.identifierTrailerFieldAccess(expr)
			if err != nil {
				return nil, err
			}
			expr = e
		}
	}

	return expr, nil
}

// expression_rhs_ident_trailer =  rhs_ident_trailer_function_call | rhs_ident_trailer_list_subscript
// right side of assignment allows function calls
func (p *Parser) expressionRhsIdentTrailer() (ast.Expr, error) {
	// consume ident
	token, ok := p.expect(tokenizer.Token_Identifier).(tokenizer.TokenWithData)
	if !ok {
		panic("bad")
	}

	var expr ast.Expr = ast.Var{
		Var: token.Value(),
	}

	trailerBeginTokens := []tokenizer.TokenKind{tokenizer.Token_LParen, tokenizer.Token_LBracket, tokenizer.Token_Dot}
	for p.nextTokenIs(trailerBeginTokens...) {
		// (1) function call
		if p.nextTokenIs(tokenizer.Token_LParen) {
			e, err := p.identifierTrailerFunctionCall(expr)
			if err != nil {
				return nil, err
			}
			expr = e
		}

		// array subscript
		if p.nextTokenIs(tokenizer.Token_LBracket) {
			e, err := p.identifierTrailerListSubscript(expr)
			if err != nil {
				return nil, err
			}
			expr = e
		}

		// object field access
		if p.nextTokenIs(tokenizer.Token_Dot) {
			e, err := p.identifierTrailerFieldAccess(expr)
			if err != nil {
				return nil, err
			}
			expr = e
		}
	}

	return expr, nil
}

// rhs_ident_trailer_function_call  := identifier { '('  expression_list `)` }?
func (p *Parser) identifierTrailerFunctionCall(expr ast.Expr) (ast.Expr, error) {
	p.expect(tokenizer.Token_LParen)
	expressionList := []ast.Expr{}
	if !p.nextTokenIs(tokenizer.Token_RParen) {
		e, err := p.expressionList()
		if err != nil {
			return nil, err
		}
		expressionList = e
	}
	p.expect(tokenizer.Token_RParen)

	return ast.FunctionCall{Expr: expr, Args: expressionList}, nil
}

// rhs_ident_trailer_list_subscript  := identifier { '['  expression`]` }?
func (p *Parser) identifierTrailerListSubscript(expr ast.Expr) (ast.Expr, error) {
	p.expect(tokenizer.Token_LBracket)
	subscriptExpr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.expect(tokenizer.Token_RBracket)

	return ast.SubscriptGet{
		Expr:      expr,
		Subscript: subscriptExpr,
	}, nil
}

// rhs_ident_trailer_field_access := identifier {'.' identifier}?
func (p *Parser) identifierTrailerFieldAccess(expr ast.Expr) (ast.Expr, error) {
	p.expect(tokenizer.Token_Dot)
	ident := p.expect(tokenizer.Token_Identifier).(tokenizer.TokenWithData)

	return ast.ObjectFieldGet{
		Expr:  expr,
		Field: ident.Value(),
	}, nil
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
