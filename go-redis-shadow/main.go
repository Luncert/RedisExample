package main

import (
	"os"
	"os/signal"
)

func main() {
	server := NewServer("localhost:7379")
	go server.Start()
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
