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
	fmt.Println(tokens)

	parser, err := parser.NewParser(tokens)
	if err != nil {
		panic(err)
	}

	ast, err := parser.Parse()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", ast)
}
