package tokenizer

import "fmt"

type TokenKind string

const (
	Token_Int   TokenKind = "Int"
	Token_Plus            = "+"
	Token_Minus           = "-"
	Token_Mul             = "*"
	Token_Div             = "/"
)

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
	return fmt.Sprintf("Simple(%+v)", t.token)
}

func (t Simple) Kind() TokenKind {
	return t.token
}
