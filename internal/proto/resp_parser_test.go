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
			// RPUSH myList 10 20
			in: "*4\r\n$5\r\nRPUSH\r\n$6\r\nmyList\r\n$2\r\n10\r\n$2\r\n20\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 4,
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
								Token:   proto.RPUSH,
								Literal: "RPUSH",
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
						Bstring: &proto.BulkStringAST{
							Prefix: proto.PrefixSymbolAST{
								Token:   proto.BULKSTRING_TOK,
								Literal: "$",
							},
							Length: proto.Digit64Bit{
								Token:   proto.DIGIT,
								Literal: 2,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "10",
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
								Literal: 2,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "20",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// LPUSHX myList 5
			in: "*3\r\n$6\r\nLPUSHX\r\n$6\r\nmyList\r\n$1\r\n5\r\n",
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
								Literal: 6,
							},
							Text: proto.StringAST{
								Token:   proto.LPUSHX,
								Literal: "LPUSHX",
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
						Bstring: &proto.BulkStringAST{
							Prefix: proto.PrefixSymbolAST{
								Token:   proto.BULKSTRING_TOK,
								Literal: "$",
							},
							Length: proto.Digit64Bit{
								Token:   proto.DIGIT,
								Literal: 1,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "5",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// LPOP myList
			in: "*2\r\n$4\r\nLPOP\r\n$6\r\nmyList\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 2,
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
								Literal: 4,
							},
							Text: proto.StringAST{
								Token:   proto.LPOP,
								Literal: "LPOP",
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
				},
			},
			err: nil,
		},
		{
			// RPOP myList
			in: "*2\r\n$4\r\nRPOP\r\n$6\r\nmyList\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 2,
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
								Literal: 4,
							},
							Text: proto.StringAST{
								Token:   proto.RPOP,
								Literal: "RPOP",
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
				},
			},
			err: nil,
		},
		{
			// BLPOP myList 5
			in: "*3\r\n$5\r\nBLPOP\r\n$6\r\nmyList\r\n$1\r\n5\r\n",
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
								Token:   proto.BLPOP,
								Literal: "BLPOP",
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
						Bstring: &proto.BulkStringAST{
							Prefix: proto.PrefixSymbolAST{
								Token:   proto.BULKSTRING_TOK,
								Literal: "$",
							},
							Length: proto.Digit64Bit{
								Token:   proto.DIGIT,
								Literal: 1,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "5",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// LLEN myList
			in: "*2\r\n$4\r\nLLEN\r\n$6\r\nmyList\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 2,
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
								Literal: 4,
							},
							Text: proto.StringAST{
								Token:   proto.LLEN,
								Literal: "LLEN",
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
				},
			},
			err: nil,
		},
		{
			// SADD mySet a b
			in: "*4\r\n$4\r\nSADD\r\n$5\r\nmySet\r\n$1\r\na\r\n$1\r\nb\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 4,
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
								Literal: 4,
							},
							Text: proto.StringAST{
								Token:   proto.SADD,
								Literal: "SADD",
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
								Literal: 5,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "mySet",
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
								Literal: 1,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "a",
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
								Literal: 1,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "b",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// SREM mySet a
			in: "*3\r\n$4\r\nSREM\r\n$5\r\nmySet\r\n$1\r\na\r\n",
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
								Literal: 4,
							},
							Text: proto.StringAST{
								Token:   proto.SREM,
								Literal: "SREM",
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
								Literal: 5,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "mySet",
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
								Literal: 1,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "a",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// SMEMBERS mySet
			in: "*2\r\n$8\r\nSMEMBERS\r\n$5\r\nmySet\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 2,
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
								Literal: 8,
							},
							Text: proto.StringAST{
								Token:   proto.SMEMBERS,
								Literal: "SMEMBERS",
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
								Literal: 5,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "mySet",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// SISMEMBER mySet a
			in: "*3\r\n$9\r\nSISMEMBER\r\n$5\r\nmySet\r\n$1\r\na\r\n",
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
								Literal: 9,
							},
							Text: proto.StringAST{
								Token:   proto.SISMEMBER,
								Literal: "SISMEMBER",
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
								Literal: 5,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "mySet",
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
								Literal: 1,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "a",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// SCARD mySet
			in: "*2\r\n$5\r\nSCARD\r\n$5\r\nmySet\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 2,
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
								Token:   proto.SCARD,
								Literal: "SCARD",
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
								Literal: 5,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "mySet",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			// SPOP mySet
			in: "*2\r\n$4\r\nSPOP\r\n$5\r\nmySet\r\n",
			out: proto.ArrayAST{
				ArrLength: proto.ArrayLengthAST{
					Prefix: proto.PrefixSymbolAST{
						Token:   proto.ARRAY_TOK,
						Literal: "*",
					},
					Length: proto.Digit64Bit{
						Token:   proto.DIGIT,
						Literal: 2,
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
								Literal: 4,
							},
							Text: proto.StringAST{
								Token:   proto.SPOP,
								Literal: "SPOP",
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
								Literal: 5,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "mySet",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			in: "*3\r\n$5\r\nLPUSH\r\n$6\r\nmyList\r\n$2\r\n30",
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
						Bstring: &proto.BulkStringAST{
							Prefix: proto.PrefixSymbolAST{
								Token:   proto.BULKSTRING_TOK,
								Literal: "$",
							},
							Length: proto.Digit64Bit{
								Token:   proto.DIGIT,
								Literal: 2,
							},
							Text: proto.StringAST{
								Token:   proto.TEXT,
								Literal: "30",
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
