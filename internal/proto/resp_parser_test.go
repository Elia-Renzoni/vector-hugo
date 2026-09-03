package proto_test

import (
	"reflect"
	"testing"

	"github.com/vector-hugo/internal/proto"
)

func TestParser(t *testing.T) {
	tests := []struct {
		in  string
		out proto.ArrayAST
		err error
	}{
		{
			in: "*3\r\n$5\r\nLPUSH\r\n$6\r\nmyList\r\n:30",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 3,
					},
				},
				Values: []proto.Literal{
					proto.Literal{
						Bstring: &proto.BulkStringAST{
							Prefix: proto.PrefixSymbolAST{
								Token:   proto.BULKSTRING_TOK,
								Literal: "$",
							},
							Length: proto.Digit64Bit{
								Token:   proto.DIGIT,
								Literal: 5,
							},
							Text: proto.StringAST{
								Token:   proto.LPUSH,
								Literal: "LPUSH",
							},
						},
					},
					proto.Literal{
						Bstring: &proto.BulkStringAST{
							Prefix: proto.PrefixSymbolAST{
								Token:   proto.BULKSTRING_TOK,
								Literal: "$",
							},
							Length: proto.Digit64Bit{
								Token:   proto.DIGIT,
								Literal: 6,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "myList",
							},
						},
					},
					proto.Literal{
						Integer: &proto.IntegerAST{
							Prefix: proto.PrefixSymbolAST{
								Token:   proto.INTEGER_TOK,
								Literal: ":",
							},
							Value: proto.Digit64Bit{
								Token:   proto.DIGIT,
								Literal: 30,
							},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		lexer := proto.NewLexer([]byte(tt.in))

		symbolList, err := lexer.Scan()
		if err != nil {
			t.Fatalf("Unexpected error while scanning the buffer: %v", err)
		}

		parser := proto.NewParser(symbolList)
		ast, err := parser.Parse()
		if err != nil {
			t.Fatalf("Unexpected error while parsing the data: %v", err)
		}

		if !reflect.DeepEqual(ast, tt.out) {
			t.Errorf("\nIncorrect Result:\nGot: %+v\nExpected:   %+v", ast, tt.out)
		}
	}
}
