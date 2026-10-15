package functions

import (
	"github.com/vector-hugo/internal/memory"
	"github.com/vector-hugo/internal/proto"
)

type (
	ModFunc  func(*memory.Arena, []string) (int, error)
	FlatFunc func(*memory.Arena) (int, error)
)

func GetFunc(cmd proto.ExecutableCommand) any {
	switch cmd.CommandName {
	case "lpush":
		return ModFunc(Lpush)
	case "rpush":
		return ModFunc(Rpush)
	case "lpushx":
		return ModFunc(Lpushx)
	case "lpop":
		return Lpop
	case "rpop":
		return Rpop
	case "llen":
		return FlatFunc(Llen)
	default:
	}

	return nil
}
