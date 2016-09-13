package weblog

import (
	"fmt"
	"net"
	"net/http"

	"github.com/deis/logger/storage"
)

const (
	bindAddr = "0.0.0.0:8088"
)

// Server implements an HTTP server.
type Server struct {
	// The Listener used for incoming HTTP connections
	Listener net.Listener
	// Server may be changed before calling Start()
	Server *http.Server
	// base URL of form http://ipaddr:port with no trailing slash
	URL string
	// started defines whether the server has started or not.
	started bool
}

// New returns a new HTTP Server. The caller should call Start to start it and Close when finished
// to shut it down.
func NewServer(storageAdapter storage.Adapter) *Server {
	s := &Server{
		Listener: defaultListener(),
		Server:   &http.Server{Handler: newRouter(newRequestHandler(storageAdapter))},
	}
	return s
}

// Start starts an HTTP server.
func (s *Server) Start() {
	if s.started {
		panic("weblog: server already started")
	}
	if s.URL == "" {
		s.URL = "http://" + s.Listener.Addr().String()
	}
	go func() {
		s.Server.Serve(s.Listener)
	}()
	s.started = true
}

// Close closes the HTTP Server from listening for the inbound requests.
func (s *Server) Close() {
	s.Server.SetKeepAlivesEnabled(false)
	s.Listener.Close()
	s.started = false
}

// defaultListener provides a net.Listener on bindAddr, panicking if it cannot listen on that
// address.
func defaultListener() net.Listener {
	l, err := net.Listen("tcp", bindAddr)
	if err != nil {
		panic(fmt.Sprintf("weblog: failed to listen on %v: %v", bindAddr, err))
	}
	return l
}
