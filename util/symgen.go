package util

import "fmt"

type SymGen struct {
	Offset int
}

func NewSymGen() SymGen {
	return SymGen{
		Offset: -1,
	}
}

func (symgen *SymGen) Next() string {
	symgen.Offset++
	return fmt.Sprintf("L%d", symgen.Offset)
}
