package tokenizer

import "fmt"

type TokenKind string

const (
	Token_Int    TokenKind = "Int"
	Token_Or     TokenKind = "or"
	Token_And    TokenKind = "or"
	Token_Lt     TokenKind = "<"
	Token_Plus   TokenKind = "+"
	Token_Minus  TokenKind = "-"
	Token_Mul    TokenKind = "*"
	Token_Div    TokenKind = "/"
	Token_Assign TokenKind = "="
	Token_LParen TokenKind = "("
	Token_RParen TokenKind = ")"
	Token_LBrace TokenKind = "{"
	Token_RBrace TokenKind = "}"

	Token_Identifier TokenKind = "Ident"
	Token_Semi       TokenKind = ";"
	Token_Comma      TokenKind = ","

	Token_Print    TokenKind = "print"
	Token_Def      TokenKind = "def"
	Token_Return   TokenKind = "return"
	Token_FuncAddr TokenKind = "function-addr"

	Token_If    TokenKind = "if"
	Token_Else  TokenKind = "else"
	Token_While TokenKind = "while"
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
}

var operators = map[string]TokenKind{
	"+": Token_Plus,
	"-": Token_Minus,
	"*": Token_Mul,
	"/": Token_Div,
	"=": Token_Assign,
	"&": Token_FuncAddr,
	"<": Token_Lt,
}

var braces = map[string]TokenKind{
	"(": Token_LParen,
	")": Token_RParen,
	"{": Token_LBrace,
	"}": Token_RBrace,
}

func SpecializeIdentifier(identifier string) Token {
	if tok, ok := specialIdentifiers[identifier]; ok {
		return Simple{
			token: tok,
		}
	}

	return TokenWithData{
		token: Token_Identifier,
		value: identifier,
	}
}

func IdentifyOperator(operator string) Token {
	if tok, ok := operators[operator]; ok {
		return &Simple{
			token: tok,
		}
	}

	return nil
}

func IdentifyBrace(operator string) Token {
	if tok, ok := braces[operator]; ok {
		return &Simple{
			token: tok,
		}
	}

	return nil
}

type Token interface {
	Kind() TokenKind
}

type TokenWithData struct {
	token TokenKind
	value string
}

func (t TokenWithData) String() string {
	return fmt.Sprintf("TokenWithData(%+v, %+v)", t.token, t.value)
}

func (t TokenWithData) Kind() TokenKind {
	return t.token
}

func (t TokenWithData) Value() string {
	return t.value
}

type Simple struct {
	token TokenKind
}

func (t Simple) String() string {
	return fmt.Sprintf("Simple('%+v')", t.token)
}

func (t Simple) Kind() TokenKind {
	return t.token
}
