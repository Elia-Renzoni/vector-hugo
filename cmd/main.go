package main

import (
	"flag"
	"log"
	"os"

	c "github.com/vector-hugo/internal/conf"
	"github.com/vector-hugo/internal/server"
)

func main() {
	path := flag.String("path", "", "a string")

	conf := c.NewConf()
	err := conf.ReadConf(*path)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	server := server.NewServer(
		conf.GetServerAddr(),
		conf.GetServerPort(),
		conf.GetConnTimeoutTime(),
		conf.GetConnMaxDataBytes(),
		conf.GetDBOptions(),
	)

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
