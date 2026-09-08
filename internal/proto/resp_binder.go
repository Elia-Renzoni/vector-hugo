package proto

import "errors"

var (
	ErrBinderExpectedBstring = errors.New("Expected a bulk string")
	ErrBinderExpectedCommand = errors.New("Expected a RESP Command")
	ErrBinderArityNotEqual   = errors.New("Arity Not Equal")
)

type ExecutableCommand struct {
	CommandType string
	Args        []string
}

var cmdArity = map[int]int{
	LPUSH:  2,
	RPUSH:  2,
	LPUSHX: 2,
	LPOP:   1,
}

func BindAST(ast ArrayAST) (ExecutableCommand, error) {
	var cmd = ExecutableCommand{}

	arrLength := ast.ArrLength.Length.Literal
	// read the Redis command
	if arrLength == 0 {
		// TODO-> define an error for this case
		return ExecutableCommand{}, nil
	}

	respCmd := ast.Values[0]
	if respCmd.Bstring != nil {
		return ExecutableCommand{}, ErrBinderExpectedBstring
	}

	arity, ok := cmdArity[respCmd.Bstring.Text.Token]
	if !ok {
		return ExecutableCommand{}, ErrBinderExpectedCommand
	}

	if arity != arrLength-1 {
		return ExecutableCommand{}, ErrBinderArityNotEqual
	}

	cmd.CommandType = respCmd.Bstring.Text.Literal
	for _, val := range ast.Values[1:] {
		cmd.Args = append(cmd.Args, val.Bstring.Text.Literal)
	}

	return cmd, nil
}
