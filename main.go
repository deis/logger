package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/crackcomm/nsqueue/consumer"
	"github.com/deis/logger/logs"
	"github.com/deis/logger/storage"
	"github.com/deis/logger/weblog"
)

var (
	// TODO: When semver permits us to do so, many of these flags should probably be phased out in
	// favor of just using environment variables.  Fewer avenues of configuring this component means
	// less confusion.
	storageType    = getopt("STORAGE_ADAPTER", "memory")
	numLines, _    = strconv.Atoi(getopt("NUMBER_OF_LINES", "1000"))
	nsqHost        = getopt("DEIS_NSQD_SERVICE_HOST", "")
	nsqPort        = getopt("DEIS_NSQD_SERVICE_PORT_TRANSPORT", "4150")
	maxInFlight, _ = strconv.Atoi(getopt("MAX_IN_FLIGHT", "30"))
	nsqURL         = fmt.Sprintf("%s:%s", nsqHost, nsqPort)
)

func main() {
	storageAdapter, err := storage.NewAdapter(storageType, numLines)
	if err != nil {
		log.Fatal("Error creating storage adapter:", err)
	}
	logger, err := logs.NewLogger(storageAdapter)
	if err != nil {
		log.Fatal("Error creating logger", err)
	}

	weblogServer, err := weblog.NewServer(logger)
	if err != nil {
		log.Fatal("Error creating weblog server", err)
	}

	weblogServer.Listen()
	log.Println("deis-logger running")

	consumer.Register("logs", "consume", maxInFlight, func(message *consumer.Message) {
		err := logger.WriteLog(message.Body)
		if err != nil {
			log.Printf("Unable to store message:%v\n%s", err, string(message.Body))
		}
		message.Success()
	})
	consumer.Connect(nsqURL)
	consumer.Start(true)
	fmt.Printf("Listening to NSQ@%s", nsqURL)
}

func handleWrite(msg *consumer.Message) {
	t := &time.Time{}
	t.UnmarshalBinary(msg.Body)
	fmt.Printf("Consume latency: %s\n", time.Since(*t))
	msg.Success()
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
