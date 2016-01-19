package weblog

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deis/logger/syslogish"
	"github.com/gorilla/mux"
)

// Server implements a simple HTTP server that handles GET and DELETE requests for application
// logs.  These actions are accomplished by delegating to a syslogish.Server, which will broker
// communication between its underlying storage.Adapter and this weblog server.
type Server struct {
	listening bool
	bindHost  string
	bindPort  int
	router    *mux.Router
}

// NewServer returns a pointer to a new Server instance.
func NewServer(bindHost string, bindPort int, syslogishServer *syslogish.Server) (*Server, error) {
	return &Server{
		bindHost: bindHost,
		bindPort: bindPort,
		router:   newRouter(newRequestHandler(syslogishServer)),
	}, nil
}

// Listen starts the server's main loop.
func (s *Server) Listen() {
	// Should only ever be called once
	if !s.listening {
		s.listening = true
		go s.listen()
		log.Printf("weblog server running on %s:%v", s.bindHost, s.bindPort)
	}
}

func (s *Server) listen() {
	http.Handle("/", s.router)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.bindHost, s.bindPort), nil); err != nil {
		log.Fatal("weblog server stopped", err)
	}
}
