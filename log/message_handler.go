package log

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/deis/logger/storage"
)

const (
	controllerPattern       = `^(INFO|WARN|DEBUG|ERROR)\s+(\[(\S+)\])+:(.*)`
	controllerContainerName = "deis-controller"
)

var (
	regex = regexp.MustCompile(controllerPattern)
)

func handle(rawMessage []byte, storageAdapter storage.Adapter) error {
	message := new(Message)
	if err := json.Unmarshal(rawMessage, message); err != nil {
		return err
	}
	if fromController(message) {
		storageAdapter.Write(getApplicationFromControllerMessage(message), message.Log)
	} else {
		labels := message.Kubernetes.Labels
		storageAdapter.Write(labels["app"], buildApplicationLogMessage(message))
	}
	return nil
}

func fromController(message *Message) bool {
	matched, _ := regexp.MatchString(controllerContainerName, message.Kubernetes.ContainerName)
	return matched
}

func getApplicationFromControllerMessage(message *Message) string {
	return regex.FindStringSubmatch(message.Log)[3]
}

func buildApplicationLogMessage(message *Message) string {
	body := message.Log
	podName := message.Kubernetes.PodName
	return fmt.Sprintf("%s -- %s", podName, body)
}
