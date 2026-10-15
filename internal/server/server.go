package server

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/vector-hugo/internal/functions"
	"github.com/vector-hugo/internal/proto"
	"github.com/vector-hugo/internal/router"
)

type VectorHugo struct {
	Address       net.IP
	ListenPort    int
	ConnsDeadline time.Time
	MaxPacketSize int // 16kb
	Dbs           *router.ListDB
}

func NewServer(ipAddr net.IP, listenPort int, deadline time.Time, maxPacketSize int) VectorHugo {
	return VectorHugo{
		Address:       ipAddr,
		ListenPort:    listenPort,
		ConnsDeadline: deadline,
		MaxPacketSize: maxPacketSize,
		Dbs:           router.NewListDB(),
	}
}

func (v VectorHugo) ListenAndServe() error {
	addr := net.JoinHostPort(v.Address.String(), strconv.Itoa(v.ListenPort))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go v.handleConn(conn)
	}
}

func (v VectorHugo) handleConn(conn net.Conn) {
	var (
		data       = make([]byte, 1024)
		bytesSum   = 0
		collection = bytes.NewBuffer(make([]byte, 0))
		chunker    proto.Chunks
		totalLines int
	)

	defer conn.Close()
	conn.SetDeadline(v.ConnsDeadline)

	for {
		n, err := conn.Read(data)
		if err != nil {
			return
		}

		bytesSum += n
		if bytesSum > v.MaxPacketSize {
			return
		}
		collection.Write(data[:n])

		if totalLines == 0 {
			elems, ok, err := proto.ReadArrayLen(collection.Bytes())
			if err != nil {
				return
			}
			if !ok {
				continue
			}
			totalLines = 1 + 2*elems
		}

		chunker.CountCompleteDataChunks(collection.Bytes())
		if chunker.ChunksCompleted(totalLines) {
			break
		}
	}

	lex := proto.NewLexer(collection.Bytes())
	tokens, err := lex.Scan()
	if err != nil {
		conn.Write([]byte("-ERR lexer: " + err.Error() + "\r\n"))
		return
	}

	parser := proto.NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		conn.Write([]byte("-ERR parser: " + err.Error() + "\r\n"))
		return
	}

	command, err := proto.BindAST(ast)
	if err != nil {
		conn.Write([]byte("-ERR binder: " + err.Error() + "\r\n"))
		return
	}

	execFunc := functions.GetFunc(command)

	create := command.CommandName == "lpush" || command.CommandName == "rpush"
	mem, err := v.Dbs.Route(command.Args[0], create)
	if err != nil {
		conn.Write([]byte("-ERR router: " + err.Error() + "\r\n"))
		return
	}

	if mem == nil {
		switch command.CommandName {
		case "lpop", "rpop":
			conn.Write([]byte("$-1\r\n")) // nil
		default: // llen, lpushx
			conn.Write([]byte(":0\r\n"))
		}
		return
	}

	var n int
	switch t := execFunc.(type) {
	case functions.FlatFunc:
		n, err = t(mem)
	case functions.ModFunc:
		n, err = t(mem, command.Args[1:])
	default:
		err = fmt.Errorf("unsupported function type %T", execFunc)
	}

	if err != nil {
		conn.Write([]byte("-ERR cmd: " + err.Error() + "\r\n"))
		return
	}

	if n == -1 {
		conn.Write([]byte("+OK\r\n"))
	} else {
		conn.Write([]byte(":" + strconv.Itoa(n) + "\r\n"))
	}
}
