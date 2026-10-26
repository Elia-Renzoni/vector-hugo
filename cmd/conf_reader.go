package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type Conf struct {
	// server configuration options
	ServerHost string `yaml:"host"`
	ServerPort int    `yaml:"port"`

	// connection configuration options
	ConnTimeoutTimeUnit  string `yaml:"timeout_unit"`
	ConnTimeoutTimeValue int    `yaml:"timeout_value"`
	ConnMaxDataBytes     int    `yaml:"max_bytes"`

	ipAddr net.IP
}

func NewConf() *Conf {
	return &Conf{}
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

	return nil
}
