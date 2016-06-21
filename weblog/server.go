package weblog

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deis/logger/logs"
	"github.com/gorilla/mux"
)

const (
	bindHost = "0.0.0.0"
	bindPort = 8088
)

// Server implements a simple HTTP server that handles GET and DELETE requests for application
// logs.  These actions are accomplished by delegating to a syslogish.Server, which will broker
// communication between its underlying storage.Adapter and this weblog server.
type Server struct {
	listening bool
	router    *mux.Router
	errCh     chan error
}

// NewServer returns a pointer to a new Server instance.
func NewServer(logger *logs.Logger) (*Server, error) {
	return &Server{
		router: newRouter(newRequestHandler(logger)),
	}, nil
}

// Listen starts the server's main loop.
func (s *Server) Listen() <-chan error {
	// Should only ever be called once
	if !s.listening {
		s.listening = true
		go func() {
			s.errCh <- s.listen()
		}()
		log.Printf("weblog server running on %s:%d", bindHost, bindPort)
	}
	return s.errCh
}

func (s *Server) listen() error {
	mux := http.NewServeMux()
	mux.Handle("/", s.router)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", bindHost, bindPort), mux)
}
