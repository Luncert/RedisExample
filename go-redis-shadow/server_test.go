package main

import (
	"log"
	"net"
	"testing"
)

const payload = "as\n"

func BenchmarkServer(b *testing.B) {
	serverAddr, _ := net.ResolveTCPAddr("tcp4", "localhost:7379")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := net.DialTCP("tcp", nil, serverAddr)
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, len(payload))
		conn.Write([]byte(payload))
		conn.Read(buf)
		conn.Close()
	}
	b.StopTimer()
}
