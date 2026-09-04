package proto

type ExecutableCommand struct {
	CommandType int
	Args        []any
}

var cmdArity = map[int]int{
	LPUSH:  2,
	RPUSH:  2,
	LPUSHX: 2,
	LPOP:   1,
}

func BindAST(ast ArrayAST) (ExecutableCommand, error) {
	var cmd = ExecutableCommand{}

	for _, val := range ast.Values {
	}

	return ExecutableCommand{}, nil
}
