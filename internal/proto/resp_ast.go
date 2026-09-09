package proto

import (
	"errors"
)

var (
	ErrASTInvalidSyntax          = errors.New("Syntax not recognised")
	ErrASTInvalidArrayToken      = errors.New("Invalid Array token: * not found")
	ErrASTInvalidTextLength      = errors.New("Mismatching string length")
	ErrASTInvalidBulkStringToken = errors.New("Invalid BulkString token: $ not found")
	ErrASTInvalidIntegerToken    = errors.New("Invalid Integer token: : not found")
	ErrASTInvalidLengthValue     = errors.New("Invalid Length value")
	ErrASTInvalidDigitType       = errors.New("Invalid digit type, must be a 64 bit digit")
	ErrASTTokensFinished         = errors.New("No more token to fetch")
)

type ArrayAST struct {
	ArrLength ArrayLengthAST
	Values    []Literal
}

type Literal struct {
	Bstring *BulkStringAST
}

type BulkStringAST struct {
	Prefix PrefixSymbolAST
	Length Digit64Bit
	Text   StringAST
}

type StringAST struct {
	Token   int
	Literal string
}

type ArrayLengthAST struct {
	Prefix PrefixSymbolAST
	Length Digit64Bit
}

type PrefixSymbolAST struct {
	Token   int
	Literal string
}

type Digit64Bit struct {
	Token   int
	Literal int
}
