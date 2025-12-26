package vm

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

// const float

type ConstFloat struct {
	Arg float64
}

func (i ConstFloat) String() string {
	return fmt.Sprintf("ConstFloat %+v", i.Arg)
}

func (i ConstFloat) isInstruction() {}

type PushString struct {
	Arg string
}

func (i PushString) String() string {
	return fmt.Sprintf("ConstString %+v", i.Arg)
}

func (i PushString) isInstruction() {}

// Asserts that there there's n entries on the evaluation stack before encountering Stack_Function_Base stack entry.

type AssertArgCount struct {
	ArgCount int
}

func (i AssertArgCount) String() string {
	return fmt.Sprintf("AssertArgCount %d", i.ArgCount)
}

func (i AssertArgCount) isInstruction() {}

type Label struct {
	Label  string
	Offset int
}

func NewLabel(label string) Label {
	return Label{
		Label:  label,
		Offset: -1,
	}

}

func (i Label) String() string {
	return fmt.Sprintf("label %d (%s):", i.Offset, i.Label)
}

func (i Label) isInstruction() {}

type Call struct {
	Offset int
	Label  string
}

func NewCall(label string) Call {
	return Call{
		Offset: -1,
		Label:  label,
	}
}

func (i Call) isInstruction() {}

func (i Call) String() string {
	return fmt.Sprintf("Call %d (%s)", i.Offset, i.Label)
}

type PushFunctionAddr struct {
	Offset int
	Label  string
}

func NewPushFunctionAddr(label string) PushFunctionAddr {
	return PushFunctionAddr{
		Offset: -1,
		Label:  label,
	}
}

func (i PushFunctionAddr) String() string {
	return fmt.Sprintf("PushFunctionAddress %+v (%s)", i.Offset, i.Label)
}

func (i PushFunctionAddr) isInstruction() {}

type BrIfNot struct {
	Label  string
	Offset int
}

func NewBrIfNot(label string) BrIfNot {
	return BrIfNot{
		Label:  label,
		Offset: -1,
	}
}

func (i BrIfNot) isInstruction() {}

func (i BrIfNot) String() string {
	return fmt.Sprintf("BrIfNot %+v (%s)", i.Offset, i.Label)
}

type BrIf struct {
	Label  string
	Offset int
}

func NewBrIf(label string) BrIf {
	return BrIf{
		Label:  label,
		Offset: -1,
	}
}

func (i BrIf) isInstruction() {}

func (i BrIf) String() string {
	return fmt.Sprintf("BrIf %+v (%s)", i.Offset, i.Label)
}

type Br struct {
	Label  string
	Offset int
}

func NewBr(label string) Br {
	return Br{
		Label:  label,
		Offset: -1,
	}
}

func (i Br) isInstruction() {}

func (i Br) String() string {
	return fmt.Sprintf("Br %+v (%s)", i.Offset, i.Label)
}

type ObjectFieldGet struct {
	Field string
}

func (i ObjectFieldGet) isInstruction() {}

func (i ObjectFieldGet) String() string {
	return fmt.Sprintf("ObjectFieldGet FIELD=%s", i.Field)
}

type ObjectFieldSet struct {
	Field string
}

func (i ObjectFieldSet) isInstruction() {}

func (i ObjectFieldSet) String() string {
	return fmt.Sprintf("ObjectFieldSet FIELD=%s", i.Field)
}

type ReadGlobal struct {
	Global string
}

func (i ReadGlobal) isInstruction() {}

func (i ReadGlobal) String() string {
	return fmt.Sprintf("ReadGlobal GLOBAL=%s", i.Global)
}

type SetGlobal struct {
	Global string
}

func (i SetGlobal) isInstruction() {}

func (i SetGlobal) String() string {
	return fmt.Sprintf("SetGlobal GLOBAL=%s", i.Global)
}
