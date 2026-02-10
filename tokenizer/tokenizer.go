package tokenizer

import (
	"fmt"
	"os"
)

type TokenReader struct {
	buf        []byte // file contents
	tokenStart int
	tokenEnd   int
	lineNumber int
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
		lineNumber: 1,
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

	return *c == ' ' || *c == '\n' || *c == '\t' || *c == '\r'
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

	return *c == '+' || *c == '-' || *c == '*' || *c == '/' || *c == '=' || *c == '&' || *c == '<' || *c == '>' || *c == '!' || *c == '%'
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
	if reader.buf[reader.tokenEnd] == '\n' {
		reader.lineNumber++
	}

	reader.tokenEnd += 1
}

func (reader *TokenReader) getToken() (string, int) {
	token := reader.buf[reader.tokenStart:reader.tokenEnd]
	reader.tokenStart = reader.tokenEnd
	return string(token), reader.lineNumber
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
		tok, lineNumber := t.reader.getToken()
		if tok == ";" {
			return NewToken(Token_Semi, lineNumber), nil
		}
		if tok == "," {
			return NewToken(Token_Comma, lineNumber), nil
		}
		if tok == "." {
			return NewToken(Token_Dot, lineNumber), nil
		}
		return nil, fmt.Errorf("Unknown punctuation")
	}

	if t.reader.peekParen() {
		t.reader.advance()
		tok, lineNumber := t.reader.getToken()

		brace, err := IdentifyBrace(tok, lineNumber)
		if err != nil {
			return nil, err
		}

		return brace, nil
	}

	if t.reader.peekOperator() {
		return t.ReadOperator()
	}

	if t.reader.peekAlpha() {
		return t.readIdentifier()
	}

	if t.reader.peekDigit() {
		return t.ReadNumber()
	}

	c := t.reader.peek()
	if *c == '"' {
		return t.readString(), nil

	}

	return nil, fmt.Errorf("Could not match any tokens")
}

func (t *Tokenizer) ReadOperator() (Token, error) {
	t.reader.advance()
	for t.reader.peekOperator() {
		t.reader.advance()
	}

	tok, lineNumber := t.reader.getToken()
	operator, err := IdentifyOperator(tok, lineNumber)
	if err != nil {
		return nil, err
	}

	return operator, nil
}

func (t *Tokenizer) ReadNumber() (Token, error) {
	t.reader.advance() // consume first digit

	for t.reader.peekDigit() {
		t.reader.advance()
	}

	tok, lineNum := t.reader.getToken()

	if *t.reader.peek() == '.' {
		t.reader.advance()

		for t.reader.peekDigit() {
			t.reader.advance()
		}

		fractionalPart, lineNum := t.reader.getToken()
		floatingPointNum := fmt.Sprintf("%s%s", tok, fractionalPart)

		return NewTokenWithValue(Token_Float, floatingPointNum, lineNum), nil
	}

	return NewTokenWithValue(Token_Int, tok, lineNum), nil
}

func (t *Tokenizer) readIdentifier() (Token, error) {
	t.reader.advance()

	for t.reader.peekAlphanumeric() {
		t.reader.advance()
	}

	tok, lineNum := t.reader.getToken()

	return SpecializeIdentifier(tok, lineNum), nil
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

// doesnt support escape sequences.
func (t Tokenizer) readString() Token {
	t.reader.advance() // consume "

	for {
		c := t.reader.peek()
		if c == nil {
			panic("reached end of file while reading the string")
		}

		if *c == '"' {
			t.reader.advance()
			strTok, lineNum := t.reader.getToken()

			trimmedString := strTok[1 : len(strTok)-1]
			return NewTokenWithValue(Token_String, trimmedString, lineNum)
		}

		t.reader.advance()
	}
}
