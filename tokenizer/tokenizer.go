package tokenizer

import (
	"fmt"
	"os"
)

type TokenReader struct {
	buf        []byte // file contents
	tokenStart int
	tokenEnd   int
}

func NewTokenReader(filename string) (*TokenReader, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &TokenReader{
		buf:        buf,
		tokenStart: 0,
		tokenEnd:   0,
	}, err
}

func (reader *TokenReader) peek() *byte {
	if reader.tokenEnd >= len(reader.buf) {
		return nil
	}

	c := reader.buf[reader.tokenEnd]
	return &c
}

func (reader *TokenReader) peekDigit() bool {
	c := reader.peek()
	if c == nil {
		return false
	}
	return *c >= '0' && *c <= '9'
}

func (reader *TokenReader) peekWhitespace() bool {
	c := reader.peek()
	if c == nil {
		return false
	}

	return *c == ' ' || *c == '\n' || *c == '\t'
}

func (reader *TokenReader) peekParen() bool {
	c := reader.peek()
	if c == nil {
		return false
	}

	return *c == '(' || *c == ')' || *c == '{' || *c == '}' || *c == '[' || *c == ']'
}

func (reader *TokenReader) peekOperator() bool {
	c := reader.peek()
	if c == nil {
		return false
	}

	return *c == '+' || *c == '-' || *c == '*' || *c == '/' || *c == '=' || *c == '&' || *c == '<' || *c == '>' || *c == '!'
}

func (reader *TokenReader) peekAlpha() bool {
	c := reader.peek()
	if c == nil {
		return false
	}

	return (*c >= 'A' && *c <= 'Z') || (*c >= 'a' && *c <= 'z') || *c == '_'
}

func (reader *TokenReader) peekPunctuation() bool {
	c := reader.peek()
	if c == nil {
		return false
	}

	return *c == ';' || *c == ',' || *c == '.'
}

func (reader *TokenReader) peekAlphanumeric() bool {
	return reader.peekAlpha() || reader.peekDigit()
}

func (reader *TokenReader) reachedEof() bool {
	return reader.tokenEnd >= len(reader.buf)
}

func (reader *TokenReader) advance() {
	reader.tokenEnd += 1
}

func (reader *TokenReader) getToken() string {
	token := reader.buf[reader.tokenStart:reader.tokenEnd]
	reader.tokenStart = reader.tokenEnd
	return string(token)
}

type Tokenizer struct {
	reader *TokenReader
}

func NewTokenizer(filename string) (*Tokenizer, error) {
	reader, err := NewTokenReader(filename)
	if err != nil {
		return nil, err
	}

	return &Tokenizer{
		reader: reader,
	}, nil
}

func (t *Tokenizer) consumeWhitespace() {
	for t.reader.peekWhitespace() {
		t.reader.advance()
	}

	t.reader.getToken()
}

func (t *Tokenizer) Tokenize() ([]Token, error) {

	tokens := make([]Token, 0)

	for true {
		tok, err := t.NextToken()
		if err != nil {
			return nil, err
		}

		if tok == nil {
			return tokens, nil
		}

		tokens = append(tokens, tok)
	}

	panic("unreachable")
}

func (t *Tokenizer) NextToken() (Token, error) {
	t.consumeWhitespaceOrComment()

	if t.reader.reachedEof() {
		return nil, nil
	}

	if t.reader.peekPunctuation() {
		t.reader.advance()
		tok := t.reader.getToken()
		if tok == ";" {
			return Simple{token: Token_Semi}, nil
		}
		if tok == "," {
			return Simple{token: Token_Comma}, nil
		}
		if tok == "." {
			return Simple{token: Token_Dot}, nil
		}

		return nil, fmt.Errorf("Unknown punctuation")
	}

	if t.reader.peekParen() {
		t.reader.advance()
		tok := t.reader.getToken()

		brace := IdentifyBrace(tok)
		if brace == nil {
			return nil, fmt.Errorf("Unknown token %+v", tok)
		}

		return brace, nil
	}

	if t.reader.peekOperator() {
		return t.ReadOperator(), nil
	}

	if t.reader.peekAlpha() {
		ident := t.readIdentifier()
		return ident, nil
	}

	if t.reader.peekDigit() {
		num := t.ReadNumber()
		return num, nil
	}

	return nil, fmt.Errorf("Could not match any tokens")
}

func (t *Tokenizer) ReadOperator() Token {
	t.reader.advance()
	for t.reader.peekOperator() {
		t.reader.advance()
	}

	tok := t.reader.getToken()
	operator := IdentifyOperator(tok)
	if operator == nil {
		panic(fmt.Sprintf("Unknown operator %+v", tok))
	}

	return operator
}

func (t *Tokenizer) ReadNumber() Token {
	t.reader.advance() // consume first digit

	for t.reader.peekDigit() {
		t.reader.advance()
	}

	tok := t.reader.getToken()

	if *t.reader.peek() == '.' {
		t.reader.advance()

		for t.reader.peekDigit() {
			t.reader.advance()
		}

		fractionalPart := t.reader.getToken()
		floatingPointNum := fmt.Sprintf("%s%s", tok, fractionalPart)

		return TokenWithData{
			token: Token_Float,
			value: floatingPointNum,
		}
	}

	return TokenWithData{
		token: Token_Int,
		value: tok,
	}
}

func (t *Tokenizer) readIdentifier() Token {
	t.reader.advance()

	for t.reader.peekAlphanumeric() {
		t.reader.advance()
	}

	tok := t.reader.getToken()

	return SpecializeIdentifier(tok)
}

func (t *Tokenizer) consumeWhitespaceOrComment() {
	t.consumeWhitespace()
	nextTok := t.reader.peek()
	for nextTok != nil && *nextTok == '#' {
		t.skipComment()
		t.consumeWhitespace()
		nextTok = t.reader.peek()
	}
}

func (t *Tokenizer) skipComment() {
	t.reader.advance() // consume #

	for true {
		tok := t.reader.peek()
		if tok == nil {
			//reached eof
			t.reader.getToken()
			return
		}
		if *tok == '\n' {
			t.reader.advance()
			t.reader.getToken()
			break
		}

		// otherwise we advance
		t.reader.advance()
	}

}
