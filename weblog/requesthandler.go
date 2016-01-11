package weblog

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/deis/logger/syslogish"
)

var getRegex *regexp.Regexp
var deleteRegex *regexp.Regexp

func init() {
	getRegex = regexp.MustCompile(`^/([-a-z0-9]+)/?(?:\?log_lines=([0-9]+))?$`)
	deleteRegex = regexp.MustCompile(`^/([-a-z0-9]+)/?$`)
}

type requestHandler struct {
	syslogishServer *syslogish.Server
}

func (h requestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.serveGet(w, r)
	} else if r.Method == "DELETE" {
		h.serveDelete(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h requestHandler) serveGet(w http.ResponseWriter, r *http.Request) {
	match := getRegex.FindStringSubmatch(r.RequestURI)
	if match == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	app := match[1]
	logLines := parseLogLines(match[2])
	logs, err := h.syslogishServer.ReadLogs(app, logLines)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), "Could not find logs for") {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	log.Printf("Returning the last %v lines for %s", logLines, app)
	for _, line := range logs {
		fmt.Fprintf(w, "%s\n", line)
	}
}

func (h requestHandler) serveDelete(w http.ResponseWriter, r *http.Request) {
	match := deleteRegex.FindStringSubmatch(r.RequestURI)
	if match == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	app := match[1]
	if err := h.syslogishServer.DestroyLogs(app); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseLogLines(logLines string) int {
	lines, err := strconv.Atoi(logLines)
	if err != nil || lines < 1 {
		log.Printf("The number of lines to return was not specified. Defaulting to 100 lines.")
		return 100
	}
	return lines
}
