package proto_test

import (
	"reflect"
	"testing"

	"github.com/vector-hugo/internal/proto"
)

func TestScan(t *testing.T) {
	tests := []struct {
		in  string
		out proto.SymbolPairs
	}{
		{
			in: "*3\r\n$5\r\nLPUSH\r\n$6\r\nmyList\r\n:30",
			out: proto.SymbolPairs{
				Tokens:   []int{proto.ARRAY_TOK, proto.DIGIT, proto.BULKSTRING_TOK, proto.DIGIT, proto.LPUSH, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT, proto.INTEGER_TOK, proto.DIGIT},
				Literals: []string{"*", "3", "$", "5", "LPUSH", "$", "6", "myList", ":", "30"},
			},
		},
	}

	for _, tt := range tests {
		lexer := proto.NewLexer([]byte(tt.in))

		got, err := lexer.Scan()
		if err != nil {
			t.Fatalf("Unexpected error while scanning the buffer: %v", err)
		}

		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("\nIncorrect Result:\nGot: %+v\nExpected:   %+v", got, tt.out)
		}
	}
}
