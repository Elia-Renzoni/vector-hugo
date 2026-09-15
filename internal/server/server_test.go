package server_test

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

type fakeRedis struct {
	addr                      string
	sendBulkstringsWithErrors []func()
}

func setupRedisClient() redis.Options {
	return redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	}
}

func setupFakeRedisClient() fakeRedis {
	return fakeRedis{
		addr: "localhost:6379",
	}
}

func TestServer(t *testing.T) {

}
