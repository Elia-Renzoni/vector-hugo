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

func (p Parser) peek() (int, string) {
	return p.symbols.Tokens[p.pos], p.symbols.Literals[p.pos]
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

	// Phase 1. Parse the first tokens of the message
	// the first ones are the array prefix ('*') and the
	// length of the message ('*<len>')

	// parse the prefix first
	tok, lit, ok := p.next()
	if !ok {
		return ast, ErrASTTokensFinished
	}
	if tok != ARRAY_TOK {
		return ast, ErrASTInvalidArrayToken
	}

	ast.ArrLength.Prefix.Token = tok
	ast.ArrLength.Prefix.Literal = lit

	// parse the length parameter
	tok, lit, ok = p.next()
	if !ok {
		return ArrayAST{}, ErrASTTokensFinished
	}

	if tok != DIGIT {
		return ArrayAST{}, ErrASTInvalidDigitType
	}

	ast.ArrLength.Length.Token = tok
	digit, err := strconv.Atoi(lit)
	if err != nil {
		return ArrayAST{}, err
	}
	ast.ArrLength.Length.Literal = digit

	// Phase 2. Parse the remaining part of the message
	// the remaining parts are composed by a fixed number of
	// literals, the expected literals are boolean types,
	// integer types and text type as well
	literals, err := p.parseLiterals()
	if err != nil {
		return ArrayAST{}, nil
	}
	ast.Values = literals

	return ast, nil
}

func (p Parser) parseLiterals() ([]Literal, error) {
	var litList = make([]Literal, 0)

	// Phase 1. Parse the prefix symbol ('$' or ':')
	tok, _, ok := p.next()
	if !ok {
		return nil, ErrASTTokensFinished
	}

	switch tok {
	case BULKSTRING_TOK:
		bs, err := p.parseBulkString()
		if err != nil {
			return nil, err
		}

		litList = append(litList, Literal{Bstring: bs})
	case INTEGER_TOK:
		intTok, err := p.parseInteger()
		if err != nil {
			return nil, err
		}

		litList = append(litList, Literal{Integer: intTok})
	default:
		return nil, ErrASTInvalidSyntax
	}

	return p.parseLiterals()
}

func (p Parser) parseBulkString() (*BulkStringAST, error) {
	var ok bool
	tok, lit := p.peek()

	bs := &BulkStringAST{}
	bs.Prefix.Token = tok
	bs.Prefix.Literal = lit

	// Phase 2. Parse the length of the upcoming string
	tok, lit, ok = p.next()
	if !ok {
		return nil, ErrASTTokensFinished
	}

	if tok != DIGIT {
		return nil, ErrASTInvalidDigitType
	}

	digit, err := strconv.Atoi(lit)
	if err != nil {
		return nil, err
	}

	bs.Length.Token = tok
	bs.Length.Literal = digit

	// Phase 3. Parse the string text. It could be
	// a simple string, like MyList, or a well-defined
	// Redis command
	tok, lit, ok = p.next()
	if !ok {
		return nil, ErrASTTokensFinished
	}

	if digit != len(lit) {
		return nil, ErrASTInvalidLengthValue
	}

	bs.Text.Token = tok
	bs.Text.Literal = lit

	return bs, nil
}

func (p Parser) parseInteger() (*IntegerAST, error) {
	var ok bool

	// Phase 2. fetch the retained token-literal pair,
	// that contains the symbol ':' token type and literal,
	// and store it inside IntegerAST
	tok, lit := p.peek()

	iast := &IntegerAST{}
	iast.Prefix.Token = tok
	iast.Prefix.Literal = lit

	// Phase 2. Parse the integer number
	tok, lit, ok = p.next()
	if !ok {
		return nil, ErrASTTokensFinished
	}

	if tok != DIGIT {
		return nil, ErrASTInvalidDigitType
	}

	digit, err := strconv.Atoi(lit)
	if err != nil {
		return nil, err
	}

	iast.Value.Token = tok
	iast.Value.Literal = digit

	return iast, nil
}
