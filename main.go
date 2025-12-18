package main

import (
	"compiler/ast"
	"compiler/eval"
	"compiler/instruction"
	"compiler/parser"
	"compiler/tokenizer"
	"fmt"
	"os"
)

const (
	DebugTokens   string = "DEBUG=tokens"
	DebugAst             = "DEBUG=ast"
	DebugBytecode        = "DEBUG=bytecode"
)

func printInstructions(instructions []instruction.Instruction) {
	for i, instruction := range instructions {
		fmt.Printf("%d: %+v\n", i, instruction)
	}
}

func main() {
	args := os.Args
	if len(args) == 1 {
		panic("Provide file.")
	}

	filename := args[1]
	// can be either DEBUG=tokens / DEBUG=ast / DEBUG=bytecode
	debugArg := ""
	if len(args) > 2 {
		debugArg = args[2]
	}

	tokenizer, err := tokenizer.NewTokenizer(filename)
	if err != nil {
		panic(err)
	}

	tokens, err := tokenizer.Tokenize()
	if err != nil {
		panic(err)
	}

	if debugArg == DebugTokens {
		fmt.Println(tokens)
		os.Exit(1)
	}

	parser, err := parser.NewParser(tokens)
	if err != nil {
		panic(err)
	}

	a, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	if debugArg == DebugAst {
		fmt.Printf("%+v\n", a)
		os.Exit(1)
	}

	visitor := ast.NewCodegenVisitor()

	a.Accept(visitor)

	if debugArg == DebugBytecode {
		printInstructions(visitor.Program.Instructions)
		os.Exit(1)
	}

	interpreter := eval.NewInterpreter(visitor.Program)
	interpreter.Run()

	/*
		rl.InitWindow(800, 600, "blah")
		for !rl.WindowShouldClose() {

			rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)
			rl.EndDrawing()
		}

		rl.CloseWindow()
	*/

}
