package main

import (
	"os"
	"os/signal"
)

func main() {
	var server Server = NewHTTPServer("localhost:7379")
	go server.Start()
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
