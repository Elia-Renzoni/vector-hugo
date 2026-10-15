package proto

import (
	"bytes"
)

func ReadWaiterChar(data []byte) (int, error) {
	input := bytes.NewBuffer(data)
	arrayCharFound := false
	for {
		char, err := input.ReadByte()
		if err != nil {
			return 0, err
		}

		if char == ARRAY {
			arrayCharFound = true
			continue
		}

		if arrayCharFound && char != ARRAY && isDigit(char) {
			return int(char), nil
		}
	}
}

type Chunks int

func (c Chunks) CountCompleteDataChunks(data []byte) {
	n := bytes.Count(data, []byte("\r\n"))
	if n > 0 {
		c = Chunks(n)
	}
}

func (c Chunks) ChunksCompleted(totalChunksExpected int) bool {
	return c >= Chunks(totalChunksExpected)
}
