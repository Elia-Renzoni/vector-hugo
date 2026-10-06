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

func NewLexer() *Lexer {
	return &Lexer{
		inputBuffer: &bytes.Buffer{},
	}
}

func (l *Lexer) addNewBuffer(buf []byte) {
	l.inputBuffer.Reset()
	l.inputBuffer.Write(buf)
}

func (l Lexer) Scan() (SymbolPairs, error) {
	var pairs = SymbolPairs{}

	for {
		char, err := l.inputBuffer.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		}

		switch char {
		case BULKSTRING, ARRAY:
			pairs.Tokens = append(pairs.Tokens, tokenResolver[string(char)], DIGIT)
		case CR, LF:
			continue
		default:
			var latestToken int = 0
			if len(pairs.Tokens) > 0 {
				latestToken = len(pairs.Tokens) - 1
			}

			switch latestToken {
			case BULKSTRING, ARRAY:
				if !isDigit(char) {
					return SymbolPairs{}, ErrLexerInvalidLiteralType
				}

				pairs.Tokens = append(pairs.Tokens, DIGIT)
				pairs.Literals = append(pairs.Literals, string(char))
			}

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

func scanLiteral(buf []byte) (string, error) {
	for _, char := range buf {
		if !isDigit(char) && !isText(char) {
			return "", ErrLexerInvalidLiteralType
		}
	}
	return string(buf), nil
}
