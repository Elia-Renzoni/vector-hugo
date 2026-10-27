package functions

import (
	"github.com/vector-hugo/internal/memory"
	"github.com/vector-hugo/internal/proto"
)

type (
	WriteFunc     func(*memory.Arena, []string) (int, error)
	ReadFunc      func(*memory.Arena) (int, error)
	DelFunc       func(*memory.Arena) (string, error)
	SeekFunc      func(*memory.Arena, int) (string, error)
	MultiSeekFunc func(*memory.Arena, ...int) ([]string, error)
)

func GetFunc(cmd proto.ExecutableCommand) any {
	switch cmd.CommandName {
	case "lpush":
		return WriteFunc(Lpush)
	case "rpush":
		return WriteFunc(Rpush)
	case "lpushx":
		return WriteFunc(Lpushx)
	case "lpop":
		return DelFunc(Lpop)
	case "rpop":
		return DelFunc(Rpop)
	case "llen":
		return ReadFunc(Llen)
	case "lindex":
		return SeekFunc(Lindex)
	case "lrange":
		return MultiSeekFunc(Lrange)
	default:
	}

	return nil
}
