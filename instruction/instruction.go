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
