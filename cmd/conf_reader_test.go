package main_test

import (
	"os"
	"path/filepath"
	"testing"

	c "github.com/vector-hugo/cmd"
)

var confExample = `
server:
host: "127.0.0.1"
port: 6379
conn:
timeout_unit: "ss"
timeout_value: 5
max_bytes: 4000
`

func TestConfReader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(path, []byte(confExample), 0o644); err != nil {
		t.Fatal(err)
	}

	conf := c.NewConf()
	err := conf.ReadConf(path)
	if err != nil {
		t.Fatal(err.Error())
	}

	if conf.ServerHost != "127.0.0.1" {
		t.Fatalf("expected host address %s but got %s", "127.0.0.1", conf.ServerHost)
	}

	if conf.ServerPort != 6379 {
		t.Fatalf("expected port %d but got %d", 6379, conf.ServerPort)
	}

	if conf.ConnTimeoutTimeUnit != "ss" {
		t.Fatalf("expected timeout time unit %s but got %s", "ss", conf.ConnTimeoutTimeUnit)
	}

	if conf.ConnTimeoutTimeValue != 5 {
		t.Fatalf("expected timeout time value %d but got %d", 5, conf.ConnTimeoutTimeValue)
	}

	if conf.ConnMaxDataBytes != 4000 {
		t.Fatalf("expected max data bytes %d but got %d", 4000, conf.ConnMaxDataBytes)
	}
}
