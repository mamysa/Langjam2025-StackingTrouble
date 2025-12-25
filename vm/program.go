package vm

type Program struct {
	Instructions            []Instruction
	EntryInstructionAddress int
	Globals                 map[string]Value
}

func (program *Program) GetInstruction(offset int) Instruction {
	// todo bounds checking
	return program.Instructions[offset]
}
