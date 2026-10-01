package server_test

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/vector-hugo/internal/server"

	rdb "github.com/redis/go-redis/v9"
)

func setupRedisClient() *rdb.Client {
	return rdb.NewClient(&rdb.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})
}

func setupVectorHugoServer() server.VectorHugo {
	return server.VectorHugo{
		Address:       net.IPv4(127, 0, 0, 1),
		ListenPort:    6379,
		ConnsDeadline: time.Time{},
		MaxPacketSize: 16000,
	}
}

type (
	CmdModifyState func(context.Context, string, ...any) *rdb.IntCmd
	CmdFetchState  func(context.Context, string) *rdb.IntCmd
)

func TestServer(t *testing.T) {
	client := setupRedisClient()
	server := setupVectorHugoServer()

	server.ListenAndServe()

	time.Sleep(3 * time.Second)

	t.Run("TestHandleConnectionWithSuccess", func(t *testing.T) {
		tests := []struct {
			in struct {
				dbName   string
				dbValues []string
			}
			cmdMod   CmdModifyState
			cmdFetch CmdFetchState
			expOut   *rdb.IntCmd
		}{
			{
				in: struct {
					dbName   string
					dbValues []string
				}{
					dbName:   "test1",
					dbValues: []string{"tvalue1", "tvalue2"},
				},
				cmdMod: client.LPush,
				expOut: rdb.NewIntResult(0, nil),
			},
			{
				in: struct {
					dbName   string
					dbValues []string
				}{
					dbName:   "test2",
					dbValues: []string{"test"},
				},
				cmdMod: client.RPush,
				expOut: rdb.NewIntResult(0, nil),
			},
			{
				in: struct {
					dbName   string
					dbValues []string
				}{
					dbName:   "test2",
					dbValues: nil,
				},
				cmdFetch: client.LLen,
				expOut:   rdb.NewIntResult(0, nil),
			},
		}

		var got *rdb.IntCmd
		for _, tt := range tests {
			if tt.cmdFetch != nil {
				got = tt.cmdFetch(context.Background(), tt.in.dbName)
			} else if tt.cmdMod != nil {
				got = tt.cmdMod(context.Background(), tt.in.dbName, tt.in.dbValues)
			}

			if !reflect.DeepEqual(got, tt.expOut) {
				t.Fatalf("expected %v but got %v", tt.expOut, got)
			}
		}
	})
}
