package parser

import (
	"compiler/ast"
	"compiler/tokenizer"
	"compiler/util"
	"fmt"
	"slices"
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

func (p *Parser) generateErrorLine(message string, expectedTokens ...tokenizer.TokenKind) error {
	const (
		RED   string = "\033[1;31m"
		GREEN        = "\033[0;32m"
		RESET        = "\033[0;0m"
	)

	currentTokenIndex := p.currentTokenIndex
	if currentTokenIndex < 0 {
		currentTokenIndex = 0
	}

	token := p.tokens[currentTokenIndex]

	var (
		lineNumberMin = max(token.LineNumber()-1, 0)
		lineNumberMax = min(token.LineNumber()+1, p.tokens[len(p.tokens)-1].LineNumber()) // last line number.
	)

	//currentTokenIndex = p.currentTokenIndex

	// read tokens

	nextTokens := make([]tokenizer.Token, 0)
	n := 2
	for true {
		nextTokenIndex := p.currentTokenIndex + n

		if nextTokenIndex >= len(p.tokens) {
			break
		}

		nextToken := p.tokens[nextTokenIndex]
		if nextToken.LineNumber() > lineNumberMax {
			break
		}

		nextTokens = append(nextTokens, nextToken)
		n = n + 1
	}

	previousTokens := make([]tokenizer.Token, 0)
	prev := 1 // start the expected token that has not been consumed yet.

	for true {
		prevTokenIndex := p.currentTokenIndex + prev

		if prevTokenIndex < 0 {
			// reached end of the string.
			break
		}

		prevToken := p.tokens[prevTokenIndex]
		if prevToken.LineNumber() < lineNumberMin {
			break
		}

		previousTokens = append(previousTokens, prevToken)

		prev = prev - 1
	}

	slices.Reverse(previousTokens)

	errorString := ""

	for i, prevToken := range previousTokens {
		if i == len(previousTokens)-1 {
			errorString = fmt.Sprintf("%s %s%s%s", errorString, RED, prevToken.Source(), RESET)
		} else {
			errorString = fmt.Sprintf("%s %s", errorString, prevToken.Source())
		}
	}

	errorString = fmt.Sprintf("%s %s%+v%s", errorString, GREEN, expectedTokens, RESET)

	for _, nextToken := range nextTokens {
		errorString = fmt.Sprintf("%s %s", errorString, nextToken.Source())
	}

	return fmt.Errorf("Parse error on line %d: %s: %s", token.LineNumber(), message, errorString)
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

func (p *Parser) expect(expectedTokenKinds ...tokenizer.TokenKind) (tokenizer.Token, error) {
	nextTokenIndex := p.currentTokenIndex + 1
	if nextTokenIndex >= len(p.tokens) {
		return nil, fmt.Errorf("Reached end of file while expecting %+v", expectedTokenKinds)
	}

	nextToken := p.tokens[nextTokenIndex]
	for _, tokenKind := range expectedTokenKinds {

		if nextToken.Kind() == tokenKind {
			p.currentTokenIndex = nextTokenIndex
			return nextToken, nil
		}
	}

	return nil, p.generateErrorLine("unexpected token", expectedTokenKinds...)
}

func (p *Parser) ensureSymbolNotDefined(functionDefs map[string]ast.FunctionDef, globalDefs *util.OrderedMap[string, ast.GlobalDef], sym string) error {
	if _, ok := functionDefs[sym]; ok {
		return fmt.Errorf("Symbol %+v redefined", sym)
	}

	if globalDefs.Contains(sym) {
		return fmt.Errorf("Symbol %+v redefined", sym)
	}

	return nil
}

// AST := FUNCTION_DEF*
func (p *Parser) Parse() (ast.Ast, error) {

	functionDefs := make(map[string]ast.FunctionDef, 0)
	globalDefs := util.NewOrderedMap[string, ast.GlobalDef]()

	for p.hasTokens() {
		if p.nextTokenIs(tokenizer.Token_Global) {
			globalDef, err := p.parseGlobalDef()
			if err != nil {
				return ast.Ast{}, err
			}

			if err := p.ensureSymbolNotDefined(functionDefs, &globalDefs, globalDef.GlobalName); err != nil {
				return ast.Ast{}, err
			}

			globalDefs.Insert(globalDef.GlobalName, globalDef)
			continue
		}

		if p.nextTokenIs(tokenizer.Token_Def) {
			functionDef, err := p.parseFunctionDef()
			if err != nil {
				return ast.Ast{}, err
			}

			if err := p.ensureSymbolNotDefined(functionDefs, &globalDefs, functionDef.Name); err != nil {
				return ast.Ast{}, err
			}

			functionDefs[functionDef.Name] = functionDef
			continue
		}

		return ast.Ast{}, p.generateErrorLine("unable to parse either function def or global", tokenizer.Token_Def, tokenizer.Token_Global)
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
	if _, err := p.expect(tokenizer.Token_Global); err != nil {
		return ast.GlobalDef{}, err
	}

	globalName, err := p.expect(tokenizer.Token_Identifier)
	if err != nil {
		return ast.GlobalDef{}, err
	}

	if _, err := p.expect(tokenizer.Token_Assign); err != nil {
		return ast.GlobalDef{}, err
	}

	expr, err := p.expression()
	if err != nil {
		return ast.GlobalDef{}, err
	}

	return ast.GlobalDef{
		GlobalName: globalName.Value(),
		Expr:       expr,
	}, err
}

// FUNC_DEF :=  'def' IDENT '(' NON_EMPTY_FUNCTION_PARAMETER_LIST? ')' statement_block
func (p *Parser) parseFunctionDef() (ast.FunctionDef, error) {
	if _, err := p.expect(tokenizer.Token_Def); err != nil {
		return ast.FunctionDef{}, err
	}

	name, err := p.expect(tokenizer.Token_Identifier)
	if err != nil {
		return ast.FunctionDef{}, err
	}

	if _, err := p.expect(tokenizer.Token_LParen); err != nil {
		return ast.FunctionDef{}, err
	}

	parameters := make([]ast.Var, 0)

	if !p.nextTokenIs(tokenizer.Token_RParen) {
		ps, err := p.parseFunctionParameterList()
		if err != nil {
			return ast.FunctionDef{}, err
		}
		parameters = ps
	}

	if _, err := p.expect(tokenizer.Token_RParen); err != nil {
		return ast.FunctionDef{}, err
	}

	statements, err := p.parseStatementBlock()
	if err != nil {
		return ast.FunctionDef{}, err
	}

	return ast.FunctionDef{
		Name:     name.Value(),
		Params:   parameters,
		StmtList: statements,
	}, nil
}

// NON_EMPTY_FUNCTION_PARAMETER_LIST = Ident {',' Ident}*
func (p *Parser) parseFunctionParameterList() ([]ast.Var, error) {
	parameters := make([]ast.Var, 0)

	ident, err := p.expect(tokenizer.Token_Identifier)
	if err != nil {
		return nil, err
	}

	param := ast.Var{
		Var: ident.Value(),
	}

	parameters = append(parameters, param)

	for p.nextTokenIs(tokenizer.Token_Comma) {
		if _, err := p.expect(tokenizer.Token_Comma); err != nil {
			return nil, err
		}
		ident, err := p.expect(tokenizer.Token_Identifier)
		if err != nil {
			return nil, err
		}

		param := ast.Var{
			Var: ident.Value(),
		}

		parameters = append(parameters, param)

	}

	return parameters, nil
}

// STATEMENT_BLOCK = '{' {PRINT_STATEMENT | ASSIGN_STATEMENT | RETURN {EXPR}?}* '}'
func (p *Parser) parseStatementBlock() ([]ast.Statement, error) {
	statements := make([]ast.Statement, 0)
	if _, err := p.expect(tokenizer.Token_LBrace); err != nil {
		return nil, err
	}

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

	if _, err := p.expect(tokenizer.Token_RBrace); err != nil {
		return nil, err
	}
	return statements, nil
}

// statement_assert = `assert` `(` expression `)`
func (p *Parser) statementAssert() (ast.Statement, error) {
	if _, err := p.expect(tokenizer.Token_Assert); err != nil {
		return nil, err
	}

	if _, err := p.expect(tokenizer.Token_LParen); err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokenizer.Token_RParen); err != nil {
		return nil, err
	}

	if _, err := p.expect(tokenizer.Token_Semi); err != nil {
		return nil, err
	}

	return ast.AssertStmt{
		Expr: expr,
	}, nil
}

// STATEMENT_WHILE `while` EXPRESSION STATEMENT_BLOCK
func (p *Parser) statementWhile() (ast.Statement, error) {
	if _, err := p.expect(tokenizer.Token_While); err != nil {
		return nil, err
	}

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
	if _, err := p.expect(tokenizer.Token_If); err != nil {
		return nil, err
	}

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
		if _, err := p.expect(tokenizer.Token_Else); err != nil {
			return nil, err
		}

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
	if _, err := p.expect(tokenizer.Token_Return); err != nil {
		return nil, err
	}

	if p.nextTokenIs(tokenizer.Token_Semi) {
		if _, err := p.expect(tokenizer.Token_Semi); err != nil {
			return nil, err
		}

		return ast.ReturnStmt{}, nil
	}

	// otherwise we have an expression
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokenizer.Token_Semi); err != nil {
		return nil, err
	}

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
		if _, err := p.expect(tokenizer.Token_LParen); err != nil {
			return nil, err
		}

		if !p.nextTokenIs(tokenizer.Token_RParen) {
			e, err := p.expressionList()
			if err != nil {
				return nil, err
			}
			expressionList = e
		}
		if _, err := p.expect(tokenizer.Token_RParen); err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_Semi); err != nil {
			return nil, err
		}
		return ast.VoidFunctionCall{
			Expr: aexpr,
			Args: expressionList,
		}, err
	}

	//definitely assignment expression
	if p.nextTokenIs(tokenizer.Token_Assign) {
		if _, err := p.expect(tokenizer.Token_Assign); err != nil {
			return nil, err
		}

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

		if _, err := p.expect(tokenizer.Token_Semi); err != nil {
			return nil, err
		}

		return ast.AssignStmt{
			AssignmentExpr: assignmentExpr,
			Expr:           expr,
		}, nil
	}

	return nil, p.generateErrorLine("unable to parse assign or function call statement", tokenizer.Token_Assign, tokenizer.Token_LParen)
}

func (p *Parser) statementPrint() (ast.Statement, error) {
	if _, err := p.expect(tokenizer.Token_Print); err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokenizer.Token_Semi); err != nil {
		return nil, err
	}

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
		if _, err := p.expect(tokenizer.Token_If); err != nil {
			return nil, err
		}

		cond, err := p.expressionDisjunction()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_Else); err != nil {
			return nil, err
		}

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
		if _, err := p.expect(tokenizer.Token_Or); err != nil {
			return nil, err
		}

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
		if _, err := p.expect(tokenizer.Token_And); err != nil {
			return nil, err
		}

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
		if _, err := p.expect(tokenizer.Token_Not); err != nil {
			return nil, err
		}

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

	expectedTokens := []tokenizer.TokenKind{tokenizer.Token_EqEq, tokenizer.Token_NotEq, tokenizer.Token_Lt, tokenizer.Token_Gt, tokenizer.Token_GrEq}
	if p.nextTokenIs(expectedTokens...) {
		token, err := p.expect(expectedTokens...)
		if err != nil {
			return nil, err
		}

		var op ast.BinOp
		switch token.Kind() {
		case tokenizer.Token_EqEq:
			op = ast.BinOp_EqEq
		case tokenizer.Token_NotEq:
			op = ast.BinOp_NotEq
		case tokenizer.Token_Lt:
			op = ast.BinOp_Lt
		case tokenizer.Token_Gt:
			op = ast.BinOp_Gt
		case tokenizer.Token_GrEq:
			op = ast.BinOp_GrEq
		default:
			return nil, p.generateErrorLine("Unknown comparison token:", expectedTokens...)
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
		opToken, err := p.expect(addOps...)
		if err != nil {
			return nil, err
		}

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
		opToken, err := p.expect(ops...)
		if err != nil {
			return nil, err
		}

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
		_, err := p.expect(tokenizer.Token_Minus)
		if err != nil {
			return nil, err
		}

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
	if p.nextTokenIs(tokenizer.Token_String) {
		tok, err := p.expect(tokenizer.Token_String)
		if err != nil {
			return nil, err
		}

		return ast.String{
			Str: tok.Value(),
		}, nil
	}

	if p.nextTokenIs(tokenizer.Token_None) {
		if _, err := p.expect(tokenizer.Token_None); err != nil {
			return nil, err
		}
		return ast.None{}, nil
	}

	if p.nextTokenIs(tokenizer.Token_True) {
		if _, err := p.expect(tokenizer.Token_True); err != nil {
			return nil, err
		}
		return ast.Bool{
			Value: true,
		}, nil
	}

	if p.nextTokenIs(tokenizer.Token_False) {
		if _, err := p.expect(tokenizer.Token_False); err != nil {
			return nil, err
		}
		return ast.Bool{
			Value: false,
		}, nil
	}

	if p.nextTokenIs(tokenizer.Token_Len) {
		if _, err := p.expect(tokenizer.Token_Len); err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_LParen); err != nil {
			return nil, err
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_RParen); err != nil {
			return nil, err
		}

		return ast.LenExpr{Expr: expr}, nil
	}

	if p.nextTokenIs(tokenizer.Token_CastFloat) {
		if _, err := p.expect(tokenizer.Token_CastFloat); err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_LParen); err != nil {
			return nil, err
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_RParen); err != nil {
			return nil, err
		}

		return ast.CastFloat{Expr: expr}, nil
	}

	if p.nextTokenIs(tokenizer.Token_Global) {
		if _, err := p.expect(tokenizer.Token_Global); err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_LParen); err != nil {
			return nil, err
		}

		globalName, err := p.expect(tokenizer.Token_Identifier)
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_RParen); err != nil {
			return nil, err
		}

		return ast.ReadGlobal{GlobalName: globalName.Value()}, nil
	}

	if p.nextTokenIs(tokenizer.Token_LBracket) {
		if _, err := p.expect(tokenizer.Token_LBracket); err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_RBracket); err != nil {
			return nil, err
		}

		return ast.NewListExpr{}, nil
	}

	if p.nextTokenIs(tokenizer.Token_LBrace) {
		if _, err := p.expect(tokenizer.Token_LBrace); err != nil {
			return nil, err
		}

		if _, err := p.expect(tokenizer.Token_RBrace); err != nil {
			return nil, err
		}

		return ast.NewObjectExpr{}, nil
	}

	if p.nextTokenIs(tokenizer.Token_FuncAddr) {
		if _, err := p.expect(tokenizer.Token_FuncAddr); err != nil {
			return nil, err
		}

		tok, err := p.expect(tokenizer.Token_Identifier)
		if err != nil {
			return nil, err
		}

		return ast.AddressOfFunction{
			Name: tok.Value(),
		}, nil
	}

	if p.nextTokenIs(tokenizer.Token_LParen) {
		_, err := p.expect(tokenizer.Token_LParen)
		if err != nil {
			return nil, err
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(tokenizer.Token_RParen)
		if err != nil {
			return nil, err
		}

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
		token, err := p.expect(tokenizer.Token_Int)
		if err != nil {
			return nil, err
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
		token, err := p.expect(tokenizer.Token_Float)
		if err != nil {
			return nil, err
		}

		parsedFloat, err := strconv.ParseFloat(token.Value(), 64)
		if err != nil {
			return nil, err
		}

		return ast.Float{
			Float: parsedFloat,
		}, nil

	}

	return nil, p.generateErrorLine(
		"unable to parse atom",
		tokenizer.Token_String,
		tokenizer.Token_None,
		tokenizer.Token_True,
		tokenizer.Token_False,
		tokenizer.Token_Len,
		tokenizer.Token_CastFloat,
		tokenizer.Token_Global,
		tokenizer.Token_LBracket,
		tokenizer.Token_LBrace,
		tokenizer.Token_FuncAddr,
		tokenizer.Token_LParen,
		tokenizer.Token_Identifier,
		tokenizer.Token_Int,
		tokenizer.Token_Float,
	)
}

// expression_lhs_ident_trailer =   rhs_ident_trailer_list_subscript
// left side of assignment doesnt allow function calls
func (p *Parser) expressionLhsIdentTrailer() (ast.Expr, error) {
	token, err := p.expect(tokenizer.Token_Identifier)
	if err != nil {
		return nil, err
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
	token, err := p.expect(tokenizer.Token_Identifier)
	if err != nil {
		return nil, err
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
	if _, err := p.expect(tokenizer.Token_LParen); err != nil {
		return nil, err
	}

	expressionList := []ast.Expr{}
	if !p.nextTokenIs(tokenizer.Token_RParen) {
		e, err := p.expressionList()
		if err != nil {
			return nil, err
		}
		expressionList = e
	}

	if _, err := p.expect(tokenizer.Token_RParen); err != nil {
		return nil, err
	}

	return ast.FunctionCall{Expr: expr, Args: expressionList}, nil
}

// rhs_ident_trailer_list_subscript  := identifier { '['  expression`]` }?
func (p *Parser) identifierTrailerListSubscript(expr ast.Expr) (ast.Expr, error) {
	if _, err := p.expect(tokenizer.Token_LBracket); err != nil {
		return nil, err
	}

	subscriptExpr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokenizer.Token_RBracket); err != nil {
		return nil, err
	}

	return ast.SubscriptGet{
		Expr:      expr,
		Subscript: subscriptExpr,
	}, nil
}

// rhs_ident_trailer_field_access := identifier {'.' identifier}?
func (p *Parser) identifierTrailerFieldAccess(expr ast.Expr) (ast.Expr, error) {
	if _, err := p.expect(tokenizer.Token_Dot); err != nil {
		return nil, err
	}

	ident, err := p.expect(tokenizer.Token_Identifier)
	if err != nil {
		return nil, err
	}

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
		if _, err := p.expect(tokenizer.Token_Comma); err != nil {
			return nil, err
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}
