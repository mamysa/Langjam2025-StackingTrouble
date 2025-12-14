package tokenizer

import "fmt"

type TokenKind string

const (
	Token_Int    TokenKind = "Int"
	Token_Plus   TokenKind = "+"
	Token_Minus  TokenKind = "-"
	Token_Mul    TokenKind = "*"
	Token_Div    TokenKind = "/"
	Token_Assign TokenKind = "="
	Token_LParen TokenKind = "("
	Token_RParen TokenKind = "("

	Token_Identifier TokenKind = "Ident"
	Token_Semi       TokenKind = ";"

	Token_Print TokenKind = "print"
)

var specialIdentifiers = map[string]TokenKind{
	"print": Token_Print,
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
