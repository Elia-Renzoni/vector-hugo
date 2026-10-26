package main

import (
	"flag"
	"log"
	"os"

	"github.com/vector-hugo/internal/server"
)

func main() {
	path := flag.String("path", "", "a string")

	conf := NewConf()
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
	)

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
