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
			in: "*3\r\n$5\r\nlpush\r\n$6\r\nmyList\r\n$2\r\n30\r\n$4\r\n3.14\r\n",
			out: proto.SymbolPairs{
				Tokens:   []int{proto.ARRAY_TOK, proto.DIGIT, proto.BULKSTRING_TOK, proto.DIGIT, proto.LPUSH, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT},
				Literals: []string{"*", "3", "$", "5", "lpush", "$", "6", "myList", "$", "2", "30", "$", "4", "3.14"},
			},
		},
		{
			in: "*3\r\n$6\r\nlindex\r\n$6\r\nmyList\r\n$2\r\n10\r\n",
			out: proto.SymbolPairs{
				Tokens:   []int{proto.ARRAY_TOK, proto.DIGIT, proto.BULKSTRING_TOK, proto.DIGIT, proto.LINDEX, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT},
				Literals: []string{"*", "3", "$", "6", "lindex", "$", "6", "myList", "$", "2", "10"},
			},
		},
		{
			in: "*4\r\n$6\r\nlrange\r\n$6\r\nmyList\r\n$1\r\n3\r\n$2\r\n10",
			out: proto.SymbolPairs{
				Tokens:   []int{proto.ARRAY_TOK, proto.DIGIT, proto.BULKSTRING_TOK, proto.DIGIT, proto.LRANGE, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT, proto.BULKSTRING_TOK, proto.DIGIT, proto.TEXT},
				Literals: []string{"*", "4", "$", "6", "lrange", "$", "6", "myList", "$", "1", "3", "$", "2", "10"},
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
