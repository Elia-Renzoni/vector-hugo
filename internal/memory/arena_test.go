package memory

import (
	"bytes"
	"testing"
)

func TestNewArena(t *testing.T) {
	a, err := NewArena()
	if err != nil {
		t.Fatalf("NewArena() error = %v", err)
	}
	defer a.DestroyArena()

	if a.mem == nil {
		t.Fatal("arena memory is nil")
	}

	expectedSize := a.maxPages * a.pageSize
	if len(a.mem) != expectedSize {
		t.Fatalf("len(mem) = %d, want %d", len(a.mem), expectedSize)
	}

	if a.memoryUsed != 0 {
		t.Fatalf("memoryUsed = %d, want 0", a.memoryUsed)
	}
}

func TestMallocAppend(t *testing.T) {
	a := newTestArena(t)

	first := []byte("hello")

	err := a.Malloc(first, true)
	if err != nil {
		t.Fatalf("Malloc() error = %v", err)
	}

	// [5][hello]
	if a.memoryUsed != 6 {
		t.Fatalf("memoryUsed = %d, want 6", a.memoryUsed)
	}

	if a.mem[0] != 5 {
		t.Fatalf("header = %d, want 5", a.mem[0])
	}

	if !bytes.Equal(a.mem[1:6], first) {
		t.Fatalf("payload = %q, want %q", a.mem[1:6], first)
	}
}

func TestMallocAppendMultiple(t *testing.T) {
	a := newTestArena(t)

	values := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("golang"),
	}

	for _, value := range values {
		if err := a.Malloc(value, true); err != nil {
			t.Fatalf("Malloc(%q) error = %v", value, err)
		}
	}

	// [5][hello][5][world][6][golang]
	expectedUsed := 6 + 6 + 7

	if a.memoryUsed != expectedUsed {
		t.Fatalf("memoryUsed = %d, want %d", a.memoryUsed, expectedUsed)
	}

	expected := []byte{
		5, 'h', 'e', 'l', 'l', 'o',
		5, 'w', 'o', 'r', 'l', 'd',
		6, 'g', 'o', 'l', 'a', 'n', 'g',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}
}

func TestMallocPrepend(t *testing.T) {
	a := newTestArena(t)

	if err := a.Malloc([]byte("hello"), true); err != nil {
		t.Fatal(err)
	}

	if err := a.Malloc([]byte("world"), false); err != nil {
		t.Fatal(err)
	}

	// [5][world][5][hello]
	expected := []byte{
		5, 'w', 'o', 'r', 'l', 'd',
		5, 'h', 'e', 'l', 'l', 'o',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}

	if a.memoryUsed != 12 {
		t.Fatalf("memoryUsed = %d, want 12", a.memoryUsed)
	}
}

func TestMallocPrependMultiple(t *testing.T) {
	a := newTestArena(t)

	values := [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	}

	for _, value := range values {
		if err := a.Malloc(value, false); err != nil {
			t.Fatalf("Malloc(%q) error = %v", value, err)
		}
	}

	// Inserendo sempre all'inizio:
	//
	// [5][three][3][two][3][one]
	expected := []byte{
		5, 't', 'h', 'r', 'e', 'e',
		3, 't', 'w', 'o',
		3, 'o', 'n', 'e',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}
}

func TestMallocEmptyPayload(t *testing.T) {
	a := newTestArena(t)

	err := a.Malloc([]byte{}, true)
	if err != nil {
		t.Fatalf("Malloc(empty) error = %v", err)
	}

	// Un payload vuoto occupa comunque il suo header:
	// [0]
	if a.memoryUsed != 1 {
		t.Fatalf("memoryUsed = %d, want 1", a.memoryUsed)
	}

	if a.mem[0] != 0 {
		t.Fatalf("header = %d, want 0", a.mem[0])
	}
}

func TestMallocMaxPayload(t *testing.T) {
	a := newTestArena(t)

	data := bytes.Repeat([]byte{'A'}, 255)

	err := a.Malloc(data, true)
	if err != nil {
		t.Fatalf("Malloc(255 bytes) error = %v", err)
	}

	if a.memoryUsed != 256 {
		t.Fatalf("memoryUsed = %d, want 256", a.memoryUsed)
	}

	if a.mem[0] != 255 {
		t.Fatalf("header = %d, want 255", a.mem[0])
	}

	if !bytes.Equal(a.mem[1:256], data) {
		t.Fatal("payload does not match")
	}
}

func TestMallocRejectsPayloadLargerThanUint8(t *testing.T) {
	a := newTestArena(t)

	data := bytes.Repeat([]byte{'A'}, 256)

	err := a.Malloc(data, true)
	if err == nil {
		t.Fatal("Malloc(256 bytes) expected error, got nil")
	}

	if a.memoryUsed != 0 {
		t.Fatalf(
			"memoryUsed = %d after rejected allocation, want 0",
			a.memoryUsed,
		)
	}
}

func TestMallocOverflow(t *testing.T) {
	a := newTestArena(t)

	// Riempiamo quasi tutta l'arena.
	data := bytes.Repeat([]byte{'A'}, a.maxPages*a.pageSize-1)

	if err := a.Malloc(data, true); err != nil {
		t.Fatalf("first Malloc() error = %v", err)
	}

	err := a.Malloc([]byte("x"), true)
	if err != ErrArenaOverflow {
		t.Fatalf(
			"second Malloc() error = %v, want %v",
			err,
			ErrArenaOverflow,
		)
	}
}

func TestFreeLastEntry(t *testing.T) {
	a := newTestArena(t)

	if err := a.Malloc([]byte("hello"), true); err != nil {
		t.Fatal(err)
	}

	if err := a.Malloc([]byte("world"), true); err != nil {
		t.Fatal(err)
	}

	// [5][hello][5][world]
	//
	// POP "world"
	if err := a.Free(nil, true); err != nil {
		t.Fatalf("Free() error = %v", err)
	}

	// Deve rimanere:
	//
	// [5][hello]
	expected := []byte{
		5, 'h', 'e', 'l', 'l', 'o',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}

	if a.memoryUsed != 6 {
		t.Fatalf("memoryUsed = %d, want 6", a.memoryUsed)
	}
}

func TestFreeFirstEntry(t *testing.T) {
	a := newTestArena(t)

	if err := a.Malloc([]byte("hello"), true); err != nil {
		t.Fatal(err)
	}

	if err := a.Malloc([]byte("world"), true); err != nil {
		t.Fatal(err)
	}

	// [5][hello][5][world]
	//
	// elimina il primo:
	//
	// [5][world]
	if err := a.Free(nil, false); err != nil {
		t.Fatalf("Free() error = %v", err)
	}

	expected := []byte{
		5, 'w', 'o', 'r', 'l', 'd',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}

	if a.memoryUsed != 6 {
		t.Fatalf("memoryUsed = %d, want 6", a.memoryUsed)
	}
}

func TestFreeTarget(t *testing.T) {
	a := newTestArena(t)

	values := [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	}

	for _, value := range values {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	// [3][one][3][two][5][three]
	//
	// elimina "two":
	//
	// [3][one][5][three]
	err := a.Free([]byte("two"), false)
	if err != nil {
		t.Fatalf("Free(target) error = %v", err)
	}

	expected := []byte{
		3, 'o', 'n', 'e',
		5, 't', 'h', 'r', 'e', 'e',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}

	if a.memoryUsed != len(expected) {
		t.Fatalf(
			"memoryUsed = %d, want %d",
			a.memoryUsed,
			len(expected),
		)
	}
}

func TestFreeTargetFirstEntry(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	err := a.Free([]byte("one"), false)
	if err != nil {
		t.Fatalf("Free(target) error = %v", err)
	}

	expected := []byte{
		3, 't', 'w', 'o',
		5, 't', 'h', 'r', 'e', 'e',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}
}

func TestFreeTargetLastEntry(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	err := a.Free([]byte("three"), false)
	if err != nil {
		t.Fatalf("Free(target) error = %v", err)
	}

	expected := []byte{
		3, 'o', 'n', 'e',
		3, 't', 'w', 'o',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}
}

func TestFreeTargetNotFound(t *testing.T) {
	a := newTestArena(t)

	if err := a.Malloc([]byte("hello"), true); err != nil {
		t.Fatal(err)
	}

	err := a.Free([]byte("missing"), false)

	if err != ErrValueNotFound {
		t.Fatalf(
			"Free(target) error = %v, want %v",
			err,
			ErrValueNotFound,
		)
	}

	// La memoria non deve essere modificata.
	expected := []byte{
		5, 'h', 'e', 'l', 'l', 'o',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf("memory changed after failed Free()")
	}
}

func TestScanFindsFirstEntry(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	found, read, err := a.Scan([]byte("one"))
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if !found {
		t.Fatal("Scan() found = false, want true")
	}

	if read != 1 {
		t.Fatalf("Scan() read = %d, want 1", read)
	}
}

func TestScanFindsMiddleEntry(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	found, read, err := a.Scan([]byte("two"))
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if !found {
		t.Fatal("Scan() found = false, want true")
	}

	if read != 2 {
		t.Fatalf("Scan() read = %d, want 2", read)
	}
}

func TestScanFindsLastEntry(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	found, read, err := a.Scan([]byte("three"))
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if !found {
		t.Fatal("Scan() found = false, want true")
	}

	if read != 3 {
		t.Fatalf("Scan() read = %d, want 3", read)
	}
}

func TestScanNotFound(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	found, read, err := a.Scan([]byte("missing"))
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if found {
		t.Fatal("Scan() found = true, want false")
	}

	if read != 3 {
		t.Fatalf("Scan() read = %d, want 3", read)
	}
}

func TestScanAfterFree(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Free([]byte("two"), false); err != nil {
		t.Fatal(err)
	}

	found, _, err := a.Scan([]byte("two"))
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if found {
		t.Fatal("Scan() found deleted value")
	}

	found, read, err := a.Scan([]byte("three"))
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if !found {
		t.Fatal("Scan() did not find remaining value")
	}

	if read != 2 {
		t.Fatalf(
			"Scan() read = %d, want 2",
			read,
		)
	}
}

func TestMallocAfterFree(t *testing.T) {
	a := newTestArena(t)

	if err := a.Malloc([]byte("hello"), true); err != nil {
		t.Fatal(err)
	}

	if err := a.Malloc([]byte("world"), true); err != nil {
		t.Fatal(err)
	}

	// Rimuove "world".
	if err := a.Free(nil, true); err != nil {
		t.Fatal(err)
	}

	// Lo spazio deve essere riutilizzabile.
	if err := a.Malloc([]byte("again"), true); err != nil {
		t.Fatalf("Malloc() after Free() error = %v", err)
	}

	expected := []byte{
		5, 'h', 'e', 'l', 'l', 'o',
		5, 'a', 'g', 'a', 'i', 'n',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}
}

func TestFreeFirstThenMalloc(t *testing.T) {
	a := newTestArena(t)

	for _, value := range [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
	} {
		if err := a.Malloc(value, true); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Free(nil, false); err != nil {
		t.Fatal(err)
	}

	if err := a.Malloc([]byte("new"), true); err != nil {
		t.Fatal(err)
	}

	expected := []byte{
		3, 't', 'w', 'o',
		5, 't', 'h', 'r', 'e', 'e',
		3, 'n', 'e', 'w',
	}

	if !bytes.Equal(a.mem[:a.memoryUsed], expected) {
		t.Fatalf(
			"memory layout = %v, want %v",
			a.mem[:a.memoryUsed],
			expected,
		)
	}
}

func newTestArena(t *testing.T) *Arena {
	t.Helper()

	a, err := NewArena()
	if err != nil {
		t.Fatalf("NewArena() error = %v", err)
	}

	t.Cleanup(func() {
		a.DestroyArena()
	})

	return a
}
