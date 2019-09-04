package main

import (
	"io/ioutil"
	"net/http"

	"github.com/Luncert/RedisExample/redishadow/log"
	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
)

// HTTPServer ...
type HTTPServer struct {
	addr    string
	router  *mux.Router
	server  *http.Server
	storage Storage
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
	log.Info("HTTPServer started at", s.addr)
	// ListenAndServe only returns
	if err := s.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}

func (s *HTTPServer) SetStorage(storage Storage) {
	s.storage = storage
}

func (s *HTTPServer) setString(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	key, err := jsonparser.GetString(data, "key")
	value, err := jsonparser.GetString(data, "value")
	if err != nil {
		// TODO: change to Debug level
		log.Fatal(err)
	}
	s.storage.SetString(key, value)
}

func (s *HTTPServer) getString(w http.ResponseWriter, r *http.Request) {
	var err error
	if key, ok := mux.Vars(r)["key"]; ok {
		if value, ok := s.storage.GetString(key); ok {
			if _, err = w.Write([]byte(value)); err != nil {
				log.Fatal("Write response data failed", err)
			}
			return
		}
	}
	if _, err = w.Write([]byte("nil")); err != nil {
		log.Fatal("Write response data", err)
	}
}

// Stop ...
func (s *HTTPServer) Stop() {
	if err := s.server.Close(); err != nil {
		log.Fatal("Stop server failed", err)
	} else {
		log.Info("HTTPServer stopped")
	}
}
