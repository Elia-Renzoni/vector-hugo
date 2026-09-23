package functions

import "github.com/vector-hugo/internal/memory"

func Lpush(bpool *memory.Arena, data []string) error {
	for _, d := range data {
		record := []byte(d)

		err := bpool.Malloc(record, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func Rpush(bpool *memory.Arena, data []string) error {
	for _, d := range data {
		record := []byte(d)

		err := bpool.Malloc(record, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func Lpushx(bpool *memory.Arena, data []string) error {
	return Lpush(bpool, data)
}

func Lpop(bpool *memory.Arena) error {
	return bpool.Free(nil, false)
}

func Rpop(bpool *memory.Arena) error {
	return bpool.Free(nil, true)
}

func Blpop(bpool *memory.Arena) {
}

func Llen(bpool *memory.Arena) (int, error) {
	_, size, err := bpool.Scan(nil)
	return size, err
}

func Sadd(bpool *memory.Arena) {
}

func Srem(bpool *memory.Arena) {
}

func Smembers(bpool *memory.Arena) {
}

func Sismember(bpool *memory.Arena) {
}

func Scard(bpool *memory.Arena) {
}

func Spop(bpool *memory.Arena) {
}
