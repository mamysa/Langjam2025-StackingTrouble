package ast

type BinOp string

type UnaryOp string

const (
	BinOp_Add BinOp = "+"
	BinOp_Sub BinOp = "-"
	BinOp_Mul BinOp = "*"
	BinOp_Div BinOp = "/"

	UnaryOp_Neg UnaryOp = "-"
)
