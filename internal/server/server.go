package server

import (
	"bytes"
	"net"
	"time"

	"github.com/vector-hugo/internal/functions"
	"github.com/vector-hugo/internal/memory"
	"github.com/vector-hugo/internal/proto"
	"github.com/vector-hugo/internal/router"
)

type VectorHugo struct {
	address       net.IP
	listenPort    int
	connsDeadline time.Time
	maxPacketSize int // 16kb
}

func NewServer(ipAddr net.IP, listenPort int, deadline time.Time, maxPacketSize int) VectorHugo {
	return VectorHugo{
		address:       ipAddr,
		listenPort:    listenPort,
		connsDeadline: deadline,
		maxPacketSize: maxPacketSize,
	}
}

func (v VectorHugo) ListenAndServe() error {
	listener, err := net.Listen("tcp", v.address.To4().String())
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

	conn.SetDeadline(v.connsDeadline)

	for {
		n, err := conn.Read(data)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				// TODO-> log the error
			}
			return
		}

		bytesSum += n
		if bytesSum >= v.maxPacketSize {
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
	dbs := router.ListDB{}

	var mem *memory.Arena

	// find the database correct database for
	// the given data structure
	if command.CommandName != "LPUSHX" {
		mem, err = dbs.Route(command.Args[0], false)
	} else {
		mem, err = dbs.Route(command.Args[0], true)
	}

	if err != nil {
		//TODO
	}

	switch t := execFunc.(type) {
	case functions.FlatFunc:
		t(mem)
	case functions.ModFunc:
		t(mem, command.Args[1:])
	}
}
