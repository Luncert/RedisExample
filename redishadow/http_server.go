package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
)

// HTTPServer ...
type HTTPServer struct {
	addr   string
	router *mux.Router
	server *http.Server
}

// NewHTTPServer ...
func NewHTTPServer(addr string) *HTTPServer {
	r := mux.NewRouter()
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return &HTTPServer{
		addr:   addr,
		router: r,
		server: s,
	}
}

// Start ...
func (s *HTTPServer) Start() {
	s.router.HandleFunc("/storage/string/set", s.setString).Methods("POST")
	s.router.HandleFunc("/storage/string/get/{key}", s.getString).Methods("GET")
	log.Println("HTTPServer started at", s.addr)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal("Start server failed", err)
	}
}

func (s *HTTPServer) setString(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	key, err := jsonparser.GetString(data, "key")
	value, err := jsonparser.GetString(data, "value")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(key, "=", value)
}

func (s *HTTPServer) getString(w http.ResponseWriter, r *http.Request) {
	if key, ok := mux.Vars(r)["key"]; ok {
		w.Write([]byte(key))
	}
}

// Stop ...
func (s *HTTPServer) Stop() {
	if err := s.server.Close(); err != nil {
		log.Fatal("Stop server failed", err)
	} else {
		log.Println("HTTPServer stopped")
	}
}
