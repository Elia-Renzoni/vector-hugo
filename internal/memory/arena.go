package memory

import (
	"bytes"
	"errors"

	"golang.org/x/sys/unix"
)

var (
	ErrArenaOverflow = errors.New("Insufficient space for storing new data")
)

type Arena struct {
	mem        []byte
	maxPages   int
	pageSize   int
	memoryUsed int
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
	if a.memoryUsed >= a.maxPages*a.pageSize || a.memoryUsed+len(data) > a.maxPages*a.pageSize {
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
	a.mem[a.memoryUsed] = memSegment.header
	copy(a.mem[a.memoryUsed+1:], memSegment.payload)

	a.memoryUsed += len(data) + 1
	return nil
}

func (a *Arena) Free() {
}

// |<h><p><h><p><h><p>| first page
func (a *Arena) Scan(target []byte) (bool, int, error) {
	dataReaded := 0
	scanOffset := 8
	for scanOffset < a.memoryUsed {
		payloadSize := int(a.mem[scanOffset])
		payloadStart := scanOffset + 1
		payloadEnd := payloadStart + payloadSize

		data := a.mem[payloadStart:payloadEnd]
		dataReaded += 1

		if bytes.Equal(data, target) {
			return true, dataReaded, nil
		}

		scanOffset = payloadEnd
	}

	return false, dataReaded, nil
}

func (a *Arena) DestroyArena() {
	unix.Munmap(a.mem)
}
