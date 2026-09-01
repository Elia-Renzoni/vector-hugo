package proto

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrLexerInvalidLiteralType error = errors.New("Mismatching token type")
)

type Lexer struct {
	inputBuffer *bytes.Buffer
}

type SymbolPairs struct {
	Tokens   []int
	Literals []string
}

func NewLexer(buf []byte) Lexer {
	return Lexer{
		inputBuffer: bytes.NewBuffer(buf),
	}
}

func (l Lexer) Scan() (SymbolPairs, error) {
	var pairs = SymbolPairs{}

	for {
		line, err := l.inputBuffer.ReadBytes(LF)
		line = bytes.TrimRight(line, "\r\n")

		if len(line) > 0 {
			char := line[0]

			switch char {
			case BULKSTRING, ARRAY, INTEGER:
				lit, err := scanDigit(line[1:])
				if err != nil {
					return SymbolPairs{}, err
				}

				pairs.Tokens = append(pairs.Tokens, tokenResolver[string(char)], DIGIT)
				pairs.Literals = append(pairs.Literals, string(char), lit)

			case SSTRING:
				lit, err := scanText(line[1:])
				if err != nil {
					return SymbolPairs{}, err
				}

				pairs.Tokens = append(pairs.Tokens, tokenResolver[string(char)], TEXT)
				pairs.Literals = append(pairs.Literals, string(char), lit)

			default:
				lit, err := scanLiteral(line)
				if err != nil {
					return SymbolPairs{}, err
				}

				token, ok := tokenResolver[lit]
				if !ok {
					pairs.Tokens = append(pairs.Tokens, TEXT)
				} else {
					pairs.Tokens = append(pairs.Tokens, token)
				}

				pairs.Literals = append(pairs.Literals, lit)
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return SymbolPairs{}, err
		}
	}

	return pairs, nil
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isText(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func scanDigit(buf []byte) (string, error) {
	for _, char := range buf {
		if !isDigit(char) {
			return "", ErrLexerInvalidLiteralType
		}
	}
	return string(buf), nil
}

func scanText(buf []byte) (string, error) {
	for _, char := range buf {
		if !isText(char) {
			return "", ErrLexerInvalidLiteralType
		}
	}
	return string(buf), nil
}

func scanLiteral(buf []byte) (string, error) {
	for _, char := range buf {
		if !isDigit(char) && !isText(char) {
			return "", ErrLexerInvalidLiteralType
		}
	}
	return string(buf), nil
}
