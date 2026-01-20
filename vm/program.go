package vm

type Program struct {
	Instructions            []Instruction
	EntryInstructionAddress int
	Globals                 map[string]value
}

func NewProgram(instructions []Instruction, entryAddress int) *Program {
	return &Program{
		Instructions:            instructions,
		EntryInstructionAddress: entryAddress,
		Globals:                 map[string]value{},
	}
}

func (program *Program) GetInstruction(offset int) Instruction {
	// todo bounds checking
	return program.Instructions[offset]
}
