package main

import (
	"bufio"
	"github.com/Luncert/RedisExample/redishadow/log"
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
	log.Info("TCPServer started at", s.addr)
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-s.stopSignal
		if err = ln.Close(); err != nil {
			log.Fatal(err)
		}
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
	log.Info("TCPServer stopped")
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
		if _, err = conn.Write([]byte(data)); err != nil {
			log.Fatal(err)
		}
	}
ret:
	conn.Close()
}

func (s *TCPServer) SetStorage(storage Storage) {
}

// Stop ...
func (s *TCPServer) Stop() {
	s.stopSignal <- true
	<-s.waitSignal
}
