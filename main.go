package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/deis/logger/syslogish"
	"github.com/deis/logger/weblog"
)

var (
	// TODO: When semver permits us to do so, many of these flags should probably be phased out in
	// favor of just using environment variables.  Fewer avenues of configuring this component means
	// less confusion.
	logAddr     = getopt("LOGGER_ADDR", "0.0.0.0")
	logHost     = getopt("HOST", "127.0.0.1")
	logPort, _  = strconv.Atoi(getopt("LOGGER_PORT", "514"))
	webAddr     = getopt("WEB_ADDR", "0.0.0.0")
	webPort, _  = strconv.Atoi(getopt("WEB_PORT", "8088"))
	storageType = getopt("STORAGE_ADAPTER", "memory")
	drainURL    = getopt("DRAIN_URL", "")
)

func main() {
	syslogishServer, err := syslogish.NewServer(logAddr, logPort, storageType, drainURL)
	if err != nil {
		log.Fatal("Error creating syslogish server", err)
	}
	weblogServer, err := weblog.NewServer(webAddr, webPort, syslogishServer)
	if err != nil {
		log.Fatal("Error creating weblog server", err)
	}

	syslogishServer.Listen()
	weblogServer.Listen()

	log.Println("deis-logger running")

	// No cleanup is needed upon termination.  The signal to reopen log files (after hypothetical
	// logroation, for instance), if applicable, is the only signal we'll care about.  Our main loop
	// will just wait for that signal.
	reopen := make(chan os.Signal, 1)
	signal.Notify(reopen, syscall.SIGUSR1)

	for {
		<-reopen
		if err := syslogishServer.ReopenLogs(); err != nil {
			log.Fatal("Error reopening logs", err)
		}
	}
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
