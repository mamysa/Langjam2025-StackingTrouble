package main

import (
	"compiler/parser"
	"compiler/tokenizer"
	"fmt"
)

func main() {
	tokenizer, err := tokenizer.NewTokenizer("test.bla")
	if err != nil {
		panic(err)
	}

	tokens, err := tokenizer.Tokenize()

	if err != nil {
		panic(err)
	}

	parser, err := parser.NewParser(tokens)
	if err != nil {
		panic(err)
	}

	ast, err := parser.Parse()

	fmt.Printf("%+v\n", ast)
}
