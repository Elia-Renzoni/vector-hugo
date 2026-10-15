package proto

import (
	"bytes"
	"errors"
	"strconv"
)

func ReadArrayLen(data []byte) (int, bool, error) {
	if len(data) == 0 {
		return 0, false, nil
	}
	if data[0] != ARRAY {
		return 0, false, errors.New("expected array")
	}
	line, _, found := bytes.Cut(data[1:], []byte("\r\n"))
	if !found {
		return 0, false, nil
	}
	n, err := strconv.Atoi(string(line))
	if err != nil {
		return 0, false, err
	}
	return n, true, nil
}

type Chunks int

func (c *Chunks) CountCompleteDataChunks(data []byte) {
	n := bytes.Count(data, []byte("\r\n"))
	if n > 0 {
		*c = Chunks(n)
	}
}

func (c Chunks) ChunksCompleted(totalChunksExpected int) bool {
	return c >= Chunks(totalChunksExpected)
}
