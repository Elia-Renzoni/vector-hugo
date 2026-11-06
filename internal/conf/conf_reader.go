package conf

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

var (
	ErrDBNotFound      = errors.New("Database not found")
	ErrDBAlreadyExists = errors.New("Database already exists")
)

type Conf struct {
	// server configuration options
	ServerHost string `yaml:"host"`
	ServerPort int    `yaml:"port"`

	// connection configuration options
	ConnTimeoutTimeUnit  string `yaml:"timeout_unit"`
	ConnTimeoutTimeValue int    `yaml:"timeout_value"`
	ConnMaxDataBytes     int    `yaml:"max_bytes"`

	// memory configuration options
	DBOptions `yaml:"database_options"`

	ipAddr net.IP
}

type MemoryOptions struct {
	MemoryMaxBytes  int  `yaml:"max_memory"`
	PageReplacer    bool `yaml:"page_replacer"`
	LRUPoolSize     int  `yaml:"lru_pool_size"`
	MemoryThreshold int  `yaml:"memory_threshold"`
}

func NewConf() *Conf {
	return &Conf{}
}

type DBOptions map[string]*MemoryOptions

func (d DBOptions) Set(dbName string, opt *MemoryOptions) error {
	if _, ok := d[dbName]; ok {
		return ErrDBAlreadyExists
	}
	d[dbName] = opt
	return nil
}

func (d DBOptions) Load(dbName string) (*MemoryOptions, error) {
	options, ok := d[dbName]
	if !ok {
		return nil, ErrDBNotFound
	}

	return options, nil
}

func (c *Conf) ReadConf(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	return c.validateOptions()
}

func (c Conf) GetServerAddr() net.IP {
	return c.ipAddr
}

func (c Conf) GetServerPort() int {
	return c.ServerPort
}

func (c Conf) GetDBOptions() DBOptions {
	return c.DBOptions
}

func (c Conf) GetConnTimeoutTime() time.Time {
	timeVal := time.Duration(c.ConnTimeoutTimeValue)
	switch c.ConnTimeoutTimeUnit {
	case SS:
		return time.Now().Add(timeVal * time.Second)
	case MS:
		return time.Now().Add(timeVal * time.Millisecond)
	}

	return time.Time{}
}

func (c Conf) GetConnMaxDataBytes() int {
	return c.ConnMaxDataBytes
}

const (
	// deadline time units
	SS string = "ss"
	MS string = "ms"
)

func (c *Conf) validateOptions() error {
	addr, err := netip.ParseAddr(c.ServerHost)
	if err != nil {
		return err
	}

	c.ipAddr = addr.AsSlice()

	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("Invalid port range values: %d", c.ServerPort)
	}

	switch c.ConnTimeoutTimeUnit {
	case "":
		c.ConnTimeoutTimeUnit = SS
	case SS, MS:
	default:
		return fmt.Errorf("Invalid connection timeout unit value: %s", c.ConnTimeoutTimeUnit)
	}

	if c.ConnTimeoutTimeValue <= 0 {
		// set the default timeout value at 10
		c.ConnTimeoutTimeValue = 10
	}

	if c.ConnMaxDataBytes <= 0 {
		// set the default max bytes at 16KB
		c.ConnMaxDataBytes = 16_000
	}

	if len(c.DBOptions) == 0 {
		c.DBOptions = make(DBOptions)
	}

	for _, opts := range c.DBOptions {
		if opts.MemoryMaxBytes <= 0 {
			opts.MemoryMaxBytes = 20480
		}

		if opts.LRUPoolSize > 0 && !opts.PageReplacer {
			return errors.New("page_replacer must be set to true")
		} else if opts.PageReplacer && opts.LRUPoolSize <= 0 {
			opts.LRUPoolSize = 16
		}

		if opts.MemoryThreshold <= 0 {
			opts.MemoryThreshold = 20480 * 0.75
		} else if opts.MemoryThreshold > opts.MemoryMaxBytes {
			return errors.New("memory_threshold must be lower than max_memory")
		}
	}

	c.DBOptions.Set("standard", &MemoryOptions{
		MemoryMaxBytes:  20480, // 20KB
		PageReplacer:    true,
		LRUPoolSize:     16,
		MemoryThreshold: 20480 * 0.75,
	})

	return nil
}
