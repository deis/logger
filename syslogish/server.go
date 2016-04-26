package syslogish

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/deis/logger/storage"
)

const (
	bindHost          = "0.0.0.0"
	bindPort          = 1514
	queueSize         = 500
	controllerPattern = `(INFO|WARN|DEBUG|ERROR)\s+(\[(\S+)\])+:(.*)`
)

var appRegex *regexp.Regexp

func init() {
	appRegex = regexp.MustCompile(`^.* ([-_a-z0-9]+)\[[a-z0-9-_\.]+\].*`)
}

// Server implements a UDP-based "syslog-like" server.  Like syslog, as described by RFC 3164, it
// expects that each packet contains a single log message and that, conversely, log messages are
// encapsulated in their entirety by a single packet, however, no attempt is made to parse the
// messages received or validate that they conform to the specification.
type Server struct {
	conn           net.PacketConn
	listening      bool
	storageQueue   chan string
	storageAdapter storage.Adapter
}

// NewServer returns a pointer to a new Server instance.
func NewServer(storageType string, numLines int) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", bindHost, bindPort))
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	newStorageAdapter, err := storage.NewAdapter(storageType, numLines)
	if err != nil {
		return nil, fmt.Errorf("configurer: Error creating storage adapter: %v", err)
	}

	return &Server{
		conn:           c,
		storageQueue:   make(chan string, queueSize),
		storageAdapter: newStorageAdapter,
	}, nil
}

// SetStorageAdapter permits a server's underlying storage.Adapter to be reconfigured (replaced)
// at runtime.
func (s *Server) SetStorageAdapter(storageAdapter storage.Adapter) {
	s.storageAdapter = storageAdapter
}

// Listen starts the server's main loop.
func (s *Server) Listen() {
	// Should only ever be called once
	if !s.listening {
		s.listening = true
		go s.receive()
		go s.processStorage()
		log.Println("syslogish server running")
	}
}

func (s *Server) receive() {
	// Make buffer the same size as the max for a UDP packet
	buf := make([]byte, 65535)
	for {
		n, _, err := s.conn.ReadFrom(buf)
		if err != nil {
			log.Fatal("syslogish server read error", err)
		}
		message := strings.TrimSuffix(string(buf[:n]), "\n")
		select {
		case s.storageQueue <- message:
		default:
		}
	}
}

func (s *Server) processStorage() {
	for message := range s.storageQueue {
		// Strip off all the leading syslog junk and just take the JSON.
		// Drop messages that clearly do not contain any JSON, although an open curly brace is only
		// a soft indicator of JSON.  If the message does not contain JSON or is otherwise malformed,
		// it may still be dropped when parsing is attempted.
		curlyIndex := strings.Index(message, "{")
		if curlyIndex > -1 {
			message = message[curlyIndex:]
			var messageJSON map[string]interface{}
			err := json.Unmarshal([]byte(message), &messageJSON)
			if err == nil && fromKubernetes(messageJSON) && !fromContainer(messageJSON, "(deis-logger.*)") {
				labels := getLabels(messageJSON)
				if fromDeisApp(labels) {
					s.storageAdapter.Write(labels["app"].(string), buildApplicationLogMessage(messageJSON))
				} else if fromController(messageJSON) {
					s.storageAdapter.Write(
						getApplicationFromControllerMessage(messageJSON),
						messageJSON["log"].(string))
				}
			}
		}
	}
}

func fromKubernetes(json map[string]interface{}) bool {
	return json["kubernetes"] != nil
}

func hasLabels(json map[string]interface{}) bool {
	return json["kubernetes"].(map[string]interface{})["labels"] != nil
}

func getLabels(json map[string]interface{}) map[string]interface{} {
	if hasLabels(json) {
		return json["kubernetes"].(map[string]interface{})["labels"].(map[string]interface{})
	}
	return nil
}

func fromDeisApp(labels map[string]interface{}) bool {
	return labels["app"] != nil &&
		labels["heritage"] != nil &&
		labels["heritage"].(string) == "deis"
}

func fromContainer(json map[string]interface{}, pattern string) bool {
	containerName := json["kubernetes"].(map[string]interface{})["container_name"].(string)
	matched, _ := regexp.MatchString(pattern, containerName)
	return matched
}

func fromController(json map[string]interface{}) bool {
	matched, _ := regexp.MatchString(controllerPattern, json["log"].(string))
	return matched
}

func getApplicationFromControllerMessage(json map[string]interface{}) string {
	regex, _ := regexp.Compile(controllerPattern)
	return regex.FindStringSubmatch(json["log"].(string))[3]
}

func buildApplicationLogMessage(json map[string]interface{}) string {
	body := json["log"].(string)
	podName := json["kubernetes"].(map[string]interface{})["pod_name"].(string)
	return fmt.Sprintf("%s -- %s", podName, body)
}

// ReadLogs returns a specified number of log lines (if available) for a specified app by
// delegating to the server's underlying storage.Adapter.
func (s *Server) ReadLogs(app string, lines int) ([]string, error) {
	if s.storageAdapter == nil {
		return nil, fmt.Errorf("Could not find logs for '%s'.  No storage adapter specified.", app)
	}
	return s.storageAdapter.Read(app, lines)
}

// DestroyLogs deletes all logs for a specified app by delegating to the server's underlying
// storage.Adapter.
func (s *Server) DestroyLogs(app string) error {
	if s.storageAdapter == nil {
		return fmt.Errorf("Could not destroy logs for '%s'.  No storage adapter specified.", app)
	}
	return s.storageAdapter.Destroy(app)
}

// ReopenLogs delegate to the server's underlying storage.Adapter to, if applicable, refresh
// references to underlying storage mechanisms.  This is useful, for instance, to ensure logging
// continues smoothly after log rotation when file-based storage is in use.
func (s *Server) ReopenLogs() error {
	if s.storageAdapter == nil {
		return errors.New("Could not reopen logs.  No storage adapter specified.")
	}
	return s.storageAdapter.Reopen()
}
