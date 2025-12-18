package eval

import "compiler/instruction"

type Program struct {
	Instructions            []instruction.Instruction
	EntryInstructionAddress int
	Globals                 map[string]Value
}

func (program *Program) GetInstruction(offset int) instruction.Instruction {
	// todo bounds checking
	return program.Instructions[offset]
}
