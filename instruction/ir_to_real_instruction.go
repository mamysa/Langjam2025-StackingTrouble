package instruction

import (
	"fmt"
)

func IrToRealInstruction(irInstructions []IrInstruction) ([]Instruction, map[string]int) {
	labelToOffsetMap := map[string]int{}

	irInstructionsWithLabelsRemoved := make([]IrInstruction, 0)

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

			irInstructionsWithLabelsRemoved = append(irInstructionsWithLabelsRemoved, nextInstruction)

			if _, ok := labelToOffsetMap[label.Label]; ok {
				panic(fmt.Errorf("label %+v is present in the map", label.Label))
			}

			// index of last instruction
			labelToOffsetMap[label.Label] = len(irInstructionsWithLabelsRemoved) - 1
			i = i + 2
			continue
		}

		// otherwise just copy the instruction over
		irInstructionsWithLabelsRemoved = append(irInstructionsWithLabelsRemoved, instruction)
		i = i + 1

	}

	// step 2: now that we have a map of label-to-offsets, replace all IrCall instances with Call with appropriate offset.

	realInstructions := make([]Instruction, 0)

	for _, irInstruction := range irInstructionsWithLabelsRemoved {
		if ircall, ok := irInstruction.(IrCall); ok {
			offset, ok := labelToOffsetMap[ircall.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", ircall.Label))
			}

			realInstructions = append(realInstructions, Call{
				Offset: offset,
			})
			continue
		}

		if pushIrFunctionAddr, ok := irInstruction.(PushIrFunctionAddr); ok {
			offset, ok := labelToOffsetMap[pushIrFunctionAddr.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", pushIrFunctionAddr.Label))
			}

			realInstructions = append(realInstructions, PushFunctionAddr{
				Offset: offset,
			})

			continue
		}

		realInstruction := irInstruction.(Instruction)
		realInstructions = append(realInstructions, realInstruction)

	}
	fmt.Println(labelToOffsetMap)

	return realInstructions, labelToOffsetMap
}
