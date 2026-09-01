package proto

import "strconv"

type Parser struct {
	symbols SymbolPairs
	pos     int
}

func NewParser(pairs SymbolPairs) Parser {
	return Parser{
		symbols: pairs,
	}
}

func (p *Parser) next() (int, string, bool) {
	if p.pos >= len(p.symbols.Tokens) {
		return 0, "", false
	}

	p.pos += 1
	return p.symbols.Tokens[p.pos], p.symbols.Literals[p.pos], true
}

func (p Parser) Parse() (ArrayAST, error) {
	var ast = ArrayAST{}

	// phase 1. fetch the array prefix first, so '*'
	tok, lit, ok := p.next()
	if !ok {
		return ast, ErrASTTokensFinished
	}
	if tok != ARRAY_TOK {
		return ast, ErrASTInvalidArrayToken
	}

	ast.arrLength.prefix.token = tok
	ast.arrLength.prefix.literal = lit

	// phase 2. fecth the actual array length, so *3\r\n the
	// method will peek the number 3
	tok, lit, ok = p.next()
	if !ok {
		return ArrayAST{}, ErrASTTokensFinished
	}

	if tok != DIGIT {
		return ArrayAST{}, ErrASTInvalidDigitType
	}

	ast.arrLength.length.token = tok
	digit, err := strconv.Atoi(lit)
	if err != nil {
		return ArrayAST{}, err
	}
	ast.arrLength.length.literal = digit

	// phase 3. fetch some bulkstring prefix symbol, for example
	// $6 fetch first '$'
	tok, lit, ok = p.next()
	if !ok {
		return ArrayAST{}, ErrASTTokensFinished
	}

	if tok != BULKSTRING_TOK {
		return ArrayAST{}, ErrASTInvalidBulkStringToken
	}

	bs := &BulkStringAST{}
	bs.prefix.token = tok
	bs.prefix.literal = lit

	// phase 4. fetch the byte size of the upcoming string
	tok, lit, ok = p.next()
	if !ok {
		return ArrayAST{}, ErrASTTokensFinished
	}

	if tok != DIGIT {
		return ArrayAST{}, ErrASTInvalidDigitType
	}

	digit, err = strconv.Atoi(lit)
	if err != nil {
		return ArrayAST{}, err
	}

	bs.length.token = tok
	bs.length.literal = digit

	// phase 5. fetch the resp command like LPUSH, LLEN, etc...
	tok, lit, ok = p.next()
	if !ok {
		return ArrayAST{}, ErrASTTokensFinished
	}

	if digit != len(lit) {
		return ArrayAST{}, ErrASTInvalidLengthValue
	}

	bs.text.token = tok
	bs.text.literal = lit

	// repeate phase 3, so fetch the next bulkstring

	return ast, nil
}
