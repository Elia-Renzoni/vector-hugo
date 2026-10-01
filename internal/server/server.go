package server

import (
	"bytes"
	"net"
	"strconv"
	"time"

	"github.com/vector-hugo/internal/functions"
	"github.com/vector-hugo/internal/memory"
	"github.com/vector-hugo/internal/proto"
	"github.com/vector-hugo/internal/router"
)

type VectorHugo struct {
	Address       net.IP
	ListenPort    int
	ConnsDeadline time.Time
	MaxPacketSize int // 16kb
	dbs           *router.ListDB
}

func NewServer(ipAddr net.IP, listenPort int, deadline time.Time, maxPacketSize int) VectorHugo {
	return VectorHugo{
		Address:       ipAddr,
		ListenPort:    listenPort,
		ConnsDeadline: deadline,
		MaxPacketSize: maxPacketSize,
		dbs:           router.NewListDB(),
	}
}

func (v VectorHugo) ListenAndServe() error {
	listener, err := net.Listen("tcp", v.Address.To4().String())
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
	)

	conn.SetDeadline(v.ConnsDeadline)

	for {
		n, err := conn.Read(data)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				// TODO-> log the error
			}
			return
		}

		bytesSum += n
		if bytesSum >= v.MaxPacketSize {
			break
		}

		collection.Write(data)
	}

	lex := proto.NewLexer(collection.Bytes())
	tokens, err := lex.Scan()
	if err != nil {
		// TODO-> handle error by sending it back to the client
		return
	}

	parser := proto.NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		// TODO-> handle error by sending it back to the client
		return
	}

	command, err := proto.BindAST(ast)
	if err != nil {
		// TODO-> handle error by sending it back to the client
		return
	}

	execFunc := functions.GetFunc(command)
	var mem *memory.Arena

	// find the database correct database for
	// the given data structure
	if command.CommandName != "LPUSHX" {
		mem, err = v.dbs.Route(command.Args[0], false)
	} else {
		mem, err = v.dbs.Route(command.Args[0], true)
	}

	if err != nil {
		//TODO
	}

	var n int = -1

	switch t := execFunc.(type) {
	case functions.FlatFunc:
		n, err = t(mem)
	case functions.ModFunc:
		err = t(mem, command.Args[1:])
	}

	if err != nil {
		conn.Write([]byte(err.Error()))
	} else {
		if n == -1 {
			conn.Write([]byte("+OK"))
		} else {
			conv := strconv.Itoa(n)
			conn.Write([]byte(conv))
		}
	}
}
