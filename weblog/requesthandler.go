package weblog

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/deis/logger/logs"
	"github.com/gorilla/mux"
)

type requestHandler struct {
	Logger *logs.Logger
}

func newRequestHandler(logger *logs.Logger) *requestHandler {
	return &requestHandler{
		Logger: logger,
	}
}

func (h requestHandler) getHealthz(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h requestHandler) getLogs(w http.ResponseWriter, req *http.Request) {
	app := mux.Vars(req)["app"]
	var logLines int
	logLinesStr := req.URL.Query().Get("log_lines")
	if logLinesStr == "" {
		log.Printf("The number of lines to return was not specified. Defaulting to 100 lines.")
		logLines = 100
	} else {
		var err error
		logLines, err = strconv.Atoi(logLinesStr)
		if err != nil {
			log.Printf("The specified number of log lines was invalid. Defaulting to 100 lines.")
			logLines = 100
		}
	}
	logs, err := h.Logger.ReadLogs(app, logLines)
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
		// strip any trailing newline characters from the logs
		fmt.Fprintf(w, "%s\n", strings.TrimSuffix(line, "\n"))
	}
}

func (h requestHandler) deleteLogs(w http.ResponseWriter, req *http.Request) {
	app := mux.Vars(req)["app"]
	if err := h.Logger.DestroyLogs(app); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
