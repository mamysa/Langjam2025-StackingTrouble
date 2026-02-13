package main

import (
	"compiler/ast"
	"compiler/parser"
	"compiler/tokenizer"
	"compiler/vm"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DebugTokens   string = "DEBUG=tokens"
	DebugAst             = "DEBUG=ast"
	DebugBytecode        = "DEBUG=bytecode"
)

func printProgram(prog *vm.Program) {
	fmt.Printf("Globals: %+v\n", prog.Globals)

	for i, instruction := range prog.Instructions {
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

	// get full path to the directory our file.bla is in. To be used for constructing paths
	// for RlLoadTexture, etc.
	fileLocation, err := filepath.Abs(filename)
	if err != nil {
		panic("unable to get absolute path for filename")
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
		printProgram(visitor.Program)
		os.Exit(1)
	}

	if err := os.Chdir(filepath.Dir(fileLocation)); err != nil {
		panic("unable to change working directory")
	}

	interpreter := vm.NewInterpreter(visitor.Program)
	interpreter.Run()
}
