package main

import (
	"github.com/Luncert/RedisExample/redishadow/log"
	"os"
	"os/signal"
)

func main() {
	log.InitLogger("./config/logger.yml")
	defer log.DestroyLogger()

	var storage Storage = NewMemoryStorage()
	var server Server = NewHTTPServer("localhost:7379")
	server.SetStorage(storage)

	go server.Start()
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
