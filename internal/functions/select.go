package functions

import (
	"strings"

	"github.com/vector-hugo/internal/memory"
	"github.com/vector-hugo/internal/proto"
)

type (
	ModFunc  func(*memory.Arena, []string) error
	FlatFunc func(*memory.Arena) (int, error)
)

func GetFunc(cmd proto.ExecutableCommand) any {
	switch strings.ToUpper(cmd.CommandName) {
	case "LPUSH":
		return Lpush
	case "RPUSH":
		return Rpush
	case "LPUSHX":
		return Lpushx
	case "LPOP":
		return Lpop
	case "RPOP":
		return Rpop
	case "LLEN":
		return Llen
	}

	return nil
}
