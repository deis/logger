package weblog

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/deis/logger/storage"
)

// TODO(bacongobbler): stop relying that port 6666 is not in use
var testBindAddr string = "127.0.0.1:6666"

// testListener provides a net.Listener for testing, panicking if it cannot listen on that
// address.
func newTestListener(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", testBindAddr)
	if err != nil {
		t.Fatalf("failed to listen on %s: %v", testBindAddr, err)
	}
	return l
}

func newTestStorageAdapter(t *testing.T) storage.Adapter {
	storageAdapter, err := storage.NewAdapter("memory", 1)
	if err != nil {
		t.Fatalf("Error creating storage adapter: %v", err)
	}
	return storageAdapter
}

func TestServerStart(t *testing.T) {
	t.Skip("skipping because of https://github.com/deis/logger/issues/120")
	storageAdapter := newTestStorageAdapter(t)
	storageAdapter.Start()
	defer storageAdapter.Stop()

	s := &Server{
		Listener: newTestListener(t),
		Server:   &http.Server{Handler: newRouter(newRequestHandler(storageAdapter))},
	}

	s.Start()

	conn, err := net.Dial("tcp", testBindAddr)
	if err != nil {
		t.Fatalf("could not connect to test server: %v", err)
	}
	defer conn.Close()
	fmt.Fprintf(conn, "GET /healthz HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Errorf("there was an error reading from the response: %v", err)
	}
	if !strings.Contains(status, "200 OK") {
		t.Errorf("Did not receive 200 OK, got '%s'", status)
	}

	// explicitly close the connection so other tests can run
	s.Close()
}

func TestServerClose(t *testing.T) {
	t.Skip("skipping because of https://github.com/deis/logger/issues/120")
	storageAdapter := newTestStorageAdapter(t)
	storageAdapter.Start()
	defer storageAdapter.Stop()

	s := &Server{
		Listener: newTestListener(t),
		Server:   &http.Server{Handler: newRouter(newRequestHandler(storageAdapter))},
	}

	s.Start()
	s.Close()

	// try reading from the server, expecting it to fail
	_, err := net.Dial("tcp", testBindAddr)
	if err == nil {
		t.Error("server returned nil. Calling s.Close() did not close the server connection!")
	}
}

func TestServerURL(t *testing.T) {
	t.Skip("skipping because of https://github.com/deis/logger/issues/120")
	storageAdapter := newTestStorageAdapter(t)
	storageAdapter.Start()
	defer storageAdapter.Stop()

	s := &Server{
		Listener: newTestListener(t),
		Server:   &http.Server{Handler: newRouter(newRequestHandler(storageAdapter))},
		URL:      "foo",
	}

	s.Start()

	if s.URL != "foo" {
		t.Errorf("URL is not 'foo', got '%s'", s.URL)
	}

	s.Close()

	s.URL = ""

	s.Start()

	if s.URL != "http://"+testBindAddr {
		t.Errorf("URL is not 'http://%s', got '%s'", testBindAddr, s.URL)
	}

	// explicitly close the connection so other tests can run
	s.Close()
}
