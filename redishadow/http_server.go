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
	log.Println("HTTPServer started at", s.addr)
	s.server.ListenAndServe()
}

func (s *HTTPServer) setString(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	data := []byte(body)
	key, err := jsonparser.GetString(data, "key")
	value, err := jsonparser.GetString(data, "value")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(key, "=", value)
}

// Stop ...
func (s *HTTPServer) Stop() {
	s.server.Close()
	log.Println("HTTPServer stoped")
}
