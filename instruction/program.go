package instruction

type Program struct {
	Instructions            []Instruction
	EntryInstructionAddress int
}

func (program *Program) GetInstruction(offset int) Instruction {
	// todo bounds checking
	return program.Instructions[offset]
}
