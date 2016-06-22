package weblog

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/deis/logger/storage"
)

const (
	bindHost = "0.0.0.0"
	bindPort = 8088
)

// Server implements a simple HTTP server that handles GET and DELETE requests for application
// logs.  These actions are accomplished by delegating to a storage.Adapter.
type Server struct {
	listening bool
	router    *mux.Router
	errCh     chan error
}

// NewServer returns a pointer to a new Server instance.
func NewServer(storageAdapter storage.Adapter) (*Server, error) {
	return &Server{
		router: newRouter(newRequestHandler(storageAdapter)),
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
