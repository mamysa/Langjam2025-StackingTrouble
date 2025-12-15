package instruction

import "fmt"

type Instruction interface {
	isInstruction()
}

// Instruction set that includes labels.
// "Real" instructions encode labels as indices in the instruction list.
type IrInstruction interface {
	isIrInstruction()
}

type InstructionNoOperands struct {
	OpCode OpCode_NoArgs
}

func (i InstructionNoOperands) isIrInstruction() {}
func (i InstructionNoOperands) isInstruction()   {}

func (i InstructionNoOperands) String() string {
	return i.OpCode.String()
}

type LoadVar struct {
	Arg string
}

func (i LoadVar) String() string {
	return fmt.Sprintf("LoadVar %s", i.Arg)
}

func (i LoadVar) isInstruction()   {}
func (i LoadVar) isIrInstruction() {}

type StoreVar struct {
	Arg string
}

func (i StoreVar) String() string {
	return fmt.Sprintf("StoreVar %s", i.Arg)
}

func (i StoreVar) isInstruction()   {}
func (i StoreVar) isIrInstruction() {}

type ConstInt struct {
	Arg int
}

func (i ConstInt) String() string {
	return fmt.Sprintf("ConstInt %d", i.Arg)
}

func (i ConstInt) isInstruction()   {}
func (i ConstInt) isIrInstruction() {}

// Asserts that there there's n entries on the evaluation stack before encountering Stack_Function_Base stack entry.

type AssertArgCount struct {
	ArgCount int
}

func (i AssertArgCount) String() string {
	return fmt.Sprintf("AssertArgCount %d", i.ArgCount)
}

func (i AssertArgCount) isInstruction()   {}
func (i AssertArgCount) isIrInstruction() {}

type Label struct {
	Label string
}

func (i Label) String() string {
	return fmt.Sprintf("label %s:", i.Label)
}

func (i Label) isIrInstruction() {}

type IrCall struct {
	Label string
}

func (i IrCall) isIrInstruction() {}

type Call struct {
	Offset int
}

func (i Call) isInstruction() {}

func (i Call) String() string {
	return fmt.Sprintf("Call %d", i.Offset)

}
