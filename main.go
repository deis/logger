package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deis/logger/syslogish"
	"github.com/deis/logger/weblog"
)

func main() {
	cfg, err := parseConfig(appName)
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	syslogishServer, err := syslogish.NewServer(cfg.StorageType, cfg.NumLines)
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
