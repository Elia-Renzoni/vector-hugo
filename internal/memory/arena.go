package memory

import (
	"bytes"
	"errors"

	"golang.org/x/sys/unix"
)

var (
	ErrArenaOverflow = errors.New("Insufficient space for storing new data")
	ErrValueNotFound = errors.New("Value not found in memory")
	ErrDataTooLarge  = errors.New("Data too large")
)

type Arena struct {
	mem               []byte
	maxPages          int
	pageSize          int
	memoryUsed        int
	latestPushedEntry memEntry
}

// |<header>|<payload>| ... |5|hello|
type memEntry struct {
	header  uint8
	payload []byte
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

func (a *Arena) Malloc(data []byte, allocStyle bool) error {
	if a.memoryUsed >= a.maxPages*a.pageSize || a.memoryUsed+len(data)+1 > a.maxPages*a.pageSize {
		return ErrArenaOverflow
	}

	if len(data) > 255 {
		return ErrDataTooLarge
	}

	memSegment := memEntry{
		header:  uint8(len(data)),
		payload: data,
	}

	// alloc the new data in an append fashion
	if allocStyle {
		// store the pairs made of an header and a related payload
		// the header indicates the length (in bytes) of the payload
		a.mem[a.memoryUsed] = memSegment.header
		copy(a.mem[a.memoryUsed+1:], memSegment.payload)

		a.memoryUsed += len(data) + 1
		a.latestPushedEntry = memSegment
		return nil
	}

	if a.memoryUsed == 0 {
		a.latestPushedEntry = memEntry{
			header:  memSegment.header,
			payload: bytes.Clone(data),
		}
	}

	// alloc the new data as the first entry
	// and shift right the oldest entries
	copy(a.mem[len(memSegment.payload)+1:], a.mem[:a.memoryUsed])
	a.mem[0] = memSegment.header
	copy(a.mem[1:], memSegment.payload)
	a.memoryUsed += len(data) + 1

	return nil
}

func (a *Arena) Free(target []byte, deletionStyle bool) ([]byte, error) {
	if target != nil {
		return a.freeWithTarget(target)
	}

	// delete the entry in the latest page, no reshape
	// nedeed
	if deletionStyle {
		data := a.latestPushedEntry.payload
		start := a.memoryUsed - len(data) - 1
		deletedData := bytes.Clone(a.mem[start+1 : a.memoryUsed])
		clear(a.mem[start:a.memoryUsed])
		a.memoryUsed = start
		return deletedData, nil
	}

	// delete the entry in the first position available,
	// in this case shifting left the remaining data is nedeed
	entrySize := 1 + int(a.mem[0])
	deletedData := bytes.Clone(a.mem[1:entrySize])
	copy(a.mem, a.mem[entrySize:a.memoryUsed])
	a.memoryUsed -= entrySize
	clear(a.mem[a.memoryUsed : a.memoryUsed+entrySize])
	return deletedData, nil
}

func (a *Arena) freeWithTarget(target []byte) ([]byte, error) {
	scanOffset := 0
	found := false
	payloadStart := 0
	payloadEnd := 0

	// Phase 1. search the target data in the buffer pool
	for scanOffset < a.memoryUsed {
		payloadSize := int(a.mem[scanOffset])
		payloadStart = scanOffset + 1
		payloadEnd = payloadStart + payloadSize

		data := a.mem[payloadStart:payloadEnd]

		if bytes.Equal(data, target) {
			found = true
			break
		}

		scanOffset = payloadEnd
	}

	if !found {
		return nil, ErrValueNotFound
	}

	// Phase 2. if founded, reshape the entries by shifting them
	// to the left after clearing the target from the pool
	entryStart := payloadStart - 1
	entryEnd := payloadEnd
	entrySize := entryEnd - entryStart

	deletedData := bytes.Clone(a.mem[payloadStart:payloadEnd])
	copy(a.mem[entryStart:], a.mem[entryEnd:a.memoryUsed])

	a.memoryUsed -= entrySize
	clear(a.mem[a.memoryUsed:])

	return deletedData, nil
}

func (a *Arena) Scan(target []byte) (bool, int, error) {
	dataReaded := 0
	scanOffset := 0
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

func (a *Arena) Seek(start, end int) ([]string, error) {
	var results []string

	dataIndex := 0
	scanOffset := 0
	for scanOffset < a.memoryUsed {
		payloadSize := int(a.mem[scanOffset])
		payloadStart := scanOffset + 1
		payloadEnd := payloadStart + payloadSize

		data := bytes.Clone(a.mem[payloadStart:payloadEnd])
		if end == 0 {
			if dataIndex == start {
				results = append(results, string(data))
				break
			}
		} else {
			if dataIndex >= start && dataIndex <= end {
				results = append(results, string(data))
			} else if dataIndex > end {
				break
			}
		}

		dataIndex += 1
		scanOffset = payloadEnd
	}

	if len(results) > 0 {
		return results, nil
	}
	return nil, ErrValueNotFound
}

func (a *Arena) DestroyArena() {
	unix.Munmap(a.mem)
}
