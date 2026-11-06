package memory_test

import (
	"reflect"
	"testing"

	"github.com/vector-hugo/internal/memory"
)

func TestLRU(t *testing.T) {
	values := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("golang"),
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	}

	lru := memory.NewLru(6)
	for _, line := range values {
		lru.Add(line)
	}

	if lru.Size() != len(values) {
		t.Fatalf("expected LRU list length of %d but got %d", len(values), lru.Size())
	}

	lru.Touch([]byte("hello"))

	got := lru.GetVictim()
	if !reflect.DeepEqual(got, []byte("world")) {
		t.Fatalf("expected %v but got %v", []byte("world"), got)
	}
}
