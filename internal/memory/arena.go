package memory

import (
	"errors"

	"golang.org/x/sys/unix"
)

var (
	ErrArenaOverflow = errors.New("Insufficient space for storing new data")
)

type Arena struct {
	mem      []byte
	maxPages int
	pageSize int
	offset   int
}

func NewArena() (*Arena, error) {
	a := &Arena{
		maxPages: 10,
		pageSize: 4096, // 4KB
	}

	// alloc the necessary memory space
	totalSpace := a.maxPages * a.pageSize
	mem, err := unix.Mmap(-1, 0, totalSpace, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANONYMOUS)
	if err != nil {
		return nil, err
	}

	a.mem = mem
	return a, nil
}

func (a *Arena) Malloc(data []byte) error {
	if a.offset >= a.maxPages*a.pageSize || a.offset+len(data) > a.maxPages*a.pageSize {
		return ErrArenaOverflow
	}

	type memEntry struct {
		header  uint8
		payload []byte
	}

	memSegment := memEntry{
		header:  uint8(len(data)),
		payload: data,
	}

	// store the pairs made of an header and a related payload
	// the header indicates the length (in bytes) of the payload
	a.mem[a.offset] = memSegment.header
	copy(a.mem[a.offset+1:], memSegment.payload)

	a.offset += len(data) + 1
	return nil
}

func (a *Arena) Free() {
}

// |<h><p><h><p><h><p>| first page
func (a *Arena) Scan(dataToSeach []byte) (int, error) {
	initialOffset := 0
	for initialOffset < a.offset {
	}

	return 0, nil
}

func (a *Arena) DestroyArena() {
	unix.Munmap(a.mem)
}
