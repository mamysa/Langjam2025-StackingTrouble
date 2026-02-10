package ast

type BinOp string

type UnaryOp string

const (
	BinOp_Add   BinOp = "+"
	BinOp_Sub   BinOp = "-"
	BinOp_Mul   BinOp = "*"
	BinOp_Div   BinOp = "/"
	BinOp_Mod   BinOp = "%"
	BinOp_Lt    BinOp = "<"
	BinOp_Gt    BinOp = ">"
	BinOp_GrEq  BinOp = ">="
	BinOp_LtEq  BinOp = "<="
	BinOp_EqEq  BinOp = "=="
	BinOp_NotEq BinOp = "!="
	BinOp_Or    BinOp = "or"
	BinOp_And   BinOp = "and"

	UnaryOp_Neg UnaryOp = "-"
	UnaryOp_Not UnaryOp = "not"
)
