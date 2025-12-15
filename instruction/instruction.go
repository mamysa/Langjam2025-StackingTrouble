package instruction

import "fmt"

type Instruction interface {
	isInstruction()
}

type InstructionNoOperands struct {
	OpCode OpCode_NoArgs
}

func (i InstructionNoOperands) isInstruction() {}

func (i InstructionNoOperands) String() string {
	return i.OpCode.String()
}

type LoadVar struct {
	Arg string
}

func (i LoadVar) String() string {
	return fmt.Sprintf("LoadVar %s", i.Arg)
}

func (i LoadVar) isInstruction() {}

type StoreVar struct {
	Arg string
}

func (i StoreVar) String() string {
	return fmt.Sprintf("StoreVar %s", i.Arg)
}

func (i StoreVar) isInstruction() {}

type ConstInt struct {
	Arg int
}

func (i ConstInt) String() string {
	return fmt.Sprintf("ConstInt %d", i.Arg)
}

func (i ConstInt) isInstruction() {}

// Asserts that there there's n entries on the evaluation stack before encountering Stack_Function_Base stack entry.

type AssertArgCount struct {
	ArgCount int
}

func (i AssertArgCount) String() string {
	return fmt.Sprintf("AssertArgCount %d", i.ArgCount)
}

func (i AssertArgCount) isInstruction() {}

type Label struct {
	Label string
}

func (i Label) String() string {
	return fmt.Sprintf("label %s:", i.Label)
}

func (i Label) isInstruction() {}
