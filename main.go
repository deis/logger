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
	storageType = getopt("STORAGE_ADAPTER", "memory")
	numLines, _ = strconv.Atoi(getopt("NUMBER_OF_LINES", "1000"))
	drainURL    = getopt("DRAIN_URL", "")
)

func main() {
	syslogishServer, err := syslogish.NewServer(storageType, numLines, drainURL)
	if err != nil {
		log.Fatal("Error creating syslogish server", err)
	}
	weblogServer, err := weblog.NewServer(syslogishServer)
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
