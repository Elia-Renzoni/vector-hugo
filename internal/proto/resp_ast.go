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
	arrLength ArrayLengthAST
	values    []Literal
}

type Literal struct {
	bstring *BulkStringAST
	integer *IntegerAST
}

type BulkStringAST struct {
	prefix PrefixSymbolAST
	length Digit64Bit
	text   StringAST
}

type StringAST struct {
	token   int
	literal string
}

type ArrayLengthAST struct {
	prefix PrefixSymbolAST
	length Digit64Bit
}

type PrefixSymbolAST struct {
	token   int
	literal string
}

type IntegerAST struct {
	prefix PrefixSymbolAST
	value  Digit64Bit
}

type Digit64Bit struct {
	token   int
	literal int
}
