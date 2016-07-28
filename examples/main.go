package main

import (
	"os"

	"github.com/pivotal-golang/lager"
	"github.com/snickers/snickers/server"
)

func main() {
	log := lager.NewLogger("snickers")
	log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	snickersServer := server.New(log, "tcp", ":8080")
	snickersServer.Start()
	select {}
}
