package conf_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	c "github.com/vector-hugo/internal/conf"
)

var confExample = `
host: 127.0.0.1
port: 6379
timeout_unit: "ss"
timeout_value: 5
max_bytes: 4000
database_options:
  myList:
    max_memory: 1024
    page_replacer: true
    lru_pool_size: 24
	memory_threshold: 512
  myList2:
    page_replacer: true
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

	if len(conf.GetDBOptions()) == 0 {
		t.Fatalf("expected options inside databsae_options")
	}

	dbName := []string{"myList", "myList2"}
	opts := conf.GetDBOptions()

	memOpts, err := opts.Load(dbName[0])
	if err != nil {
		t.Fatal(err)
	}

	expectedMemOpts := c.MemoryOptions{MemoryMaxBytes: 1024, PageReplacer: true, LRUPoolSize: 24, MemoryThreshold: 512}
	if !reflect.DeepEqual(memOpts, &expectedMemOpts) {
		t.Fatalf("expected memory options: %#v but got: %#v", expectedMemOpts, memOpts)
	}

	memOpts, err = opts.Load(dbName[1])
	if err != nil {
		t.Fatal(err)
	}

	expectedMemOpts = c.MemoryOptions{MemoryMaxBytes: 20480, PageReplacer: true, LRUPoolSize: 16, MemoryThreshold: 15360}
	if !reflect.DeepEqual(memOpts, &expectedMemOpts) {
		t.Fatalf("expected memory options: %#v but got: %#v", expectedMemOpts, memOpts)
	}

}
