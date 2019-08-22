package main

import (
	"log"
	"net"
	"testing"
)

var serverAddr *net.TCPAddr

func BenchmarkServer(b *testing.B) {
	jobEngine := NewJobEngine(128, 32)

	server := NewServer("localhost:7379")
	go server.Start()
	serverAddr, _ = net.ResolveTCPAddr("tcp4", "localhost:7379")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := jobEngine.Execute(mockClientRequest, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	jobEngine.Shutdown()
	b.StopTimer()

	server.Stop()
}

func mockClientRequest(...interface{}) {
	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 4)
	conn.Write([]byte("Hi!\n"))
	conn.Read(buf)
}
