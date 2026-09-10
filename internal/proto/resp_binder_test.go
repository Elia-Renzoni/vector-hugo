package proto_test

import (
	"reflect"
	"testing"

	"github.com/vector-hugo/internal/proto"
)

func TestBindAST(t *testing.T) {
	tests := []struct {
		in  proto.ArrayAST
		out proto.ExecutableCommand
		err error
	}{
		{
			in: proto.ArrayAST{
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
			out: proto.ExecutableCommand{
				CommandName: "LPUSH",
				Args:        []string{"myList", "30"},
			},
			err: nil,
		},
		{
			// LPOP myList
			in: proto.ArrayAST{
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
			out: proto.ExecutableCommand{
				CommandName: "LPOP",
				Args:        []string{"myList"},
			},
			err: nil,
		},
	}

	for _, tt := range tests {

		execCmd, err := proto.BindAST(tt.in)
		if err != nil {
			t.Fatalf("Unexpected error while binding the data: %v", err)
		}

		if !reflect.DeepEqual(execCmd, tt.out) {
			t.Errorf("\nIncorrect Result:\nGot: %+v\nExpected:   %+v", execCmd, tt.out)
		}
	}

}
