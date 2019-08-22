package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// Server ...
type Server struct {
	addr       string
	stopSignal chan bool
	waitSignal chan bool
}

// NewServer ...
func NewServer(addr string) *Server {
	return &Server{
		addr:       addr,
		stopSignal: make(chan bool, 1),
		waitSignal: make(chan bool, 1),
	}
}

// Start ...
func (s *Server) Start() {
	log.Println("Server started")
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
	log.Println("Server stoped")
}

func (s *Server) handleConnection(conn net.Conn) {
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
func (s *Server) Stop() {
	s.stopSignal <- true
	<-s.waitSignal
}
