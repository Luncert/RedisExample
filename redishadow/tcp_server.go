package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// TCPServer ...
type TCPServer struct {
	addr       string
	stopSignal chan bool
	waitSignal chan bool
}

// NewTCPServer ...
func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		addr:       addr,
		stopSignal: make(chan bool, 1),
		waitSignal: make(chan bool, 1),
	}
}

// Start ...
func (s *TCPServer) Start() {
	log.Println("TCPServer started at", s.addr)
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-s.stopSignal
		ln.Close()
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			errStr := err.Error()
			switch {
			case strings.Contains(errStr, "closed network connection"):
				goto ret
			default:
				log.Fatal(err)
			}
		}
		go s.handleConnection(conn)
	}
ret:
	s.waitSignal <- true
	log.Println("TCPServer stoped")
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		data, err := reader.ReadString('\n')
		if err != nil {
			errStr := err.Error()
			switch {
			case errStr == "EOF":
				goto ret
			default:
				log.Fatal(err)
			}
		}
		conn.Write([]byte(data))
	}
ret:
	conn.Close()
}

// Stop ...
func (s *TCPServer) Stop() {
	s.stopSignal <- true
	<-s.waitSignal
}
