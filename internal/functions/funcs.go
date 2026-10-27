package functions

import (
	"bytes"
	"errors"

	"github.com/vector-hugo/internal/memory"
)

var (
	ErrEmptyReturnedData = errors.New("Empty data")
	ErrListOutOfRange    = errors.New("List out of range")
)

func Lpush(bpool *memory.Arena, data []string) (int, error) {
	for _, d := range data {
		record := []byte(d)

		err := bpool.Malloc(record, false)
		if err != nil {
			return 0, err
		}
	}

	return len(data), nil
}

func Rpush(bpool *memory.Arena, data []string) (int, error) {
	for _, d := range data {
		record := []byte(d)

		err := bpool.Malloc(record, true)
		if err != nil {
			return 0, err
		}
	}

	return len(data), nil
}

func Lpushx(bpool *memory.Arena, data []string) (int, error) {
	return Lpush(bpool, data)
}

func Lpop(bpool *memory.Arena) (string, error) {
	data, err := bpool.Free(nil, false)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", ErrEmptyReturnedData
	}
	result := bytes.NewBuffer(data)
	return result.String(), err
}

func Rpop(bpool *memory.Arena) (string, error) {
	data, err := bpool.Free(nil, true)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", ErrEmptyReturnedData
	}

	result := bytes.NewBuffer(data)
	return result.String(), nil
}

func Llen(bpool *memory.Arena) (int, error) {
	_, size, err := bpool.Scan(nil)
	return size, err
}

func Lindex(bpool *memory.Arena, position int) (string, error) {
	if position < 0 {
		return "", ErrListOutOfRange
	}

	result, err := bpool.Seek(position)
	if err != nil {
		return "", err
	}

	return result[0], nil
}
