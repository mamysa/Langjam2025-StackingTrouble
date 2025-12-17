package main

import (
	"compiler/ast"
	"compiler/eval"
	"compiler/instruction"
	"compiler/parser"
	"compiler/tokenizer"
	"fmt"
)

func printInstructions(instructions []instruction.Instruction) {
	for i, instruction := range instructions {
		fmt.Printf("%d: %+v\n", i, instruction)
	}
}

func main() {
	tokenizer, err := tokenizer.NewTokenizer("test2.bla")
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

	a, err := parser.Parse()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", a)

	visitor := ast.CodegenVisitor{}

	a.Accept(&visitor)
	printInstructions(visitor.Program.Instructions)

	interpreter := eval.NewInterpreter(visitor.Program)
	interpreter.Run()
}
