package instruction

import "fmt"

func IrToRealInstruction(irInstructions []IrInstruction) ([]Instruction, map[string]int) {
	labelToOffsetMap := map[string]int{}

	irInstructionsWithLabelsRemoved := make([]Instruction, 0)

	i := 0
	for i < len(irInstructions) {
		instruction := irInstructions[i]
		label, ok := instruction.(Label)
		if ok {
			nextInstructionIndex := i + 1
			if nextInstructionIndex >= len(irInstructions) {
				panic("label not followed by any instructions")
			}
			nextInstruction := irInstructions[nextInstructionIndex]

			// two consecutive labels, crash.
			if _, ok := nextInstruction.(Label); ok {
				panic("double label")
			}

			assertRealNextInstr := nextInstruction.(Instruction)

			irInstructionsWithLabelsRemoved = append(irInstructionsWithLabelsRemoved, assertRealNextInstr)

			if _, ok := labelToOffsetMap[label.Label]; ok {
				panic(fmt.Errorf("label %+v is present in the map", label.Label))
			}

			// index of last instruction
			labelToOffsetMap[label.Label] = len(irInstructionsWithLabelsRemoved) - 1
			i = i + 2
			continue
		}

		assertRealInstr := instruction.(Instruction)

		// otherwise just copy the instruction over
		irInstructionsWithLabelsRemoved = append(irInstructionsWithLabelsRemoved, assertRealInstr)
		i = i + 1

	}

	fmt.Println(labelToOffsetMap)

	return irInstructionsWithLabelsRemoved, labelToOffsetMap
}
