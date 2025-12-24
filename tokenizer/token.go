package tokenizer

import "fmt"

type TokenKind string

const (
	Token_String TokenKind = "String"
	Token_Int    TokenKind = "Int"
	Token_Float  TokenKind = "Float"
	Token_Or     TokenKind = "or"
	Token_And    TokenKind = "and"
	Token_Not    TokenKind = "not"

	Token_Lt       TokenKind = "<"
	Token_Gt       TokenKind = ">"
	Token_GrEq     TokenKind = ">="
	Token_EqEq     TokenKind = "=="
	Token_NotEq    TokenKind = "!="
	Token_Plus     TokenKind = "+"
	Token_Minus    TokenKind = "-"
	Token_Mul      TokenKind = "*"
	Token_Div      TokenKind = "/"
	Token_Assign   TokenKind = "="
	Token_LParen   TokenKind = "("
	Token_RParen   TokenKind = ")"
	Token_LBrace   TokenKind = "{"
	Token_RBrace   TokenKind = "}"
	Token_LBracket TokenKind = "["
	Token_RBracket TokenKind = "]"

	Token_Identifier TokenKind = "Ident"
	Token_Semi       TokenKind = ";"
	Token_Comma      TokenKind = ","
	Token_Dot        TokenKind = "."

	Token_Print    TokenKind = "print"
	Token_Def      TokenKind = "def"
	Token_Return   TokenKind = "return"
	Token_FuncAddr TokenKind = "&"

	Token_If    TokenKind = "if"
	Token_Else  TokenKind = "else"
	Token_While TokenKind = "while"

	Token_Len       TokenKind = "len"
	Token_Assert    TokenKind = "assert"
	Token_CastFloat TokenKind = "float"

	Token_True  TokenKind = "True"
	Token_False TokenKind = "False"

	Token_None TokenKind = "None"

	Token_Global TokenKind = "global"
)

var specialIdentifiers = map[string]TokenKind{
	"print":  Token_Print,
	"def":    Token_Def,
	"return": Token_Return,
	"if":     Token_If,
	"else":   Token_Else,
	"while":  Token_While,
	"or":     Token_Or,
	"and":    Token_And,
	"not":    Token_Not,
	"len":    Token_Len,
	"assert": Token_Assert,
	"True":   Token_True,
	"False":  Token_False,
	"None":   Token_None,
	"global": Token_Global,
	"float":  Token_CastFloat,
}

var operators = map[string]TokenKind{
	"+":  Token_Plus,
	"-":  Token_Minus,
	"*":  Token_Mul,
	"/":  Token_Div,
	"=":  Token_Assign,
	"&":  Token_FuncAddr,
	"<":  Token_Lt,
	">":  Token_Gt,
	"==": Token_EqEq,
	"!=": Token_NotEq,
	">=": Token_GrEq,
}

var braces = map[string]TokenKind{
	"(": Token_LParen,
	")": Token_RParen,
	"{": Token_LBrace,
	"}": Token_RBrace,
	"[": Token_LBracket,
	"]": Token_RBracket,
}

func SpecializeIdentifier(identifier string, lineNumber int) Token {
	if tok, ok := specialIdentifiers[identifier]; ok {
		return NewToken(tok, lineNumber)
	}

	return NewTokenWithValue(Token_Identifier, identifier, lineNumber)
}

func IdentifyOperator(operator string, lineNumber int) (Token, error) {
	if tok, ok := operators[operator]; ok {
		return NewToken(tok, lineNumber), nil
	}

	return nil, fmt.Errorf("Error processing symbol %+v on line %+v", operator, lineNumber)
}

func IdentifyBrace(brace string, lineNumber int) (Token, error) {
	if tok, ok := braces[brace]; ok {
		return NewToken(tok, lineNumber), nil
	}

	return nil, fmt.Errorf("Error processing symbol %+v on line %+v", brace, lineNumber)
}

type Token interface {
	Kind() TokenKind
	Value() string
	LineNumber() int
	Source() string // returns symbol for the token.
}

type simpleToken struct {
	tokenKind  TokenKind
	lineNumber int
}

func NewToken(tokenKind TokenKind, lineNumber int) Token {
	return simpleToken{
		tokenKind:  tokenKind,
		lineNumber: lineNumber,
	}
}

func (t simpleToken) Kind() TokenKind {
	return t.tokenKind
}

func (t simpleToken) Value() string {
	panic("This token doesn't have attached value")
}

func (t simpleToken) LineNumber() int {
	return t.lineNumber
}

func (t simpleToken) String() string {
	return fmt.Sprintf("Token(%+v, lineNumber: %+v)", t.tokenKind, t.lineNumber)
}

func (t simpleToken) Source() string {
	return string(t.tokenKind)
}

type tokenWithValue struct {
	token      TokenKind
	value      string
	lineNumber int
}

func NewTokenWithValue(tokenKind TokenKind, value string, lineNumber int) Token {
	return tokenWithValue{
		token:      tokenKind,
		value:      value,
		lineNumber: lineNumber,
	}
}

func (t tokenWithValue) Kind() TokenKind {
	return t.token
}

func (t tokenWithValue) Value() string {
	return t.value
}

func (t tokenWithValue) LineNumber() int {
	return t.lineNumber
}

func (t tokenWithValue) Source() string {
	return t.value
}

func (t tokenWithValue) String() string {
	return fmt.Sprintf("Token(%+v, value: %+v, lineNumber: %+v)", t.token, t.value, t.lineNumber)
}
