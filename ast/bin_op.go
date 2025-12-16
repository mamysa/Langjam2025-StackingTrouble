package ast

type BinOp string

type UnaryOp string

const (
	BinOp_Add BinOp = "+"
	BinOp_Sub BinOp = "-"
	BinOp_Mul BinOp = "*"
	BinOp_Div BinOp = "/"
	BinOp_Lt  BinOp = "<"
	BinOp_Or  BinOp = "or"
	BinOp_And BinOp = "and"

	UnaryOp_Neg UnaryOp = "-"
)
