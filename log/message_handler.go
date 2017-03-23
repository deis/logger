package log

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/deis/logger/storage"
)

const (
	podPattern              = `(\w.*)-(\w.*)-(\w.*)-(\w.*)`
	controllerPattern       = `^(INFO|WARN|DEBUG|ERROR)\s+(\[(\S+)\])+:(.*)`
	controllerContainerName = "deis-controller"
	timeFormat              = "2006-01-02T15:04:05-07:00"
)

var (
	controllerRegex = regexp.MustCompile(controllerPattern)
	podRegex        = regexp.MustCompile(podPattern)
)

func handle(rawMessage []byte, storageAdapter storage.Adapter) error {
	message := new(Message)
	if err := json.Unmarshal(rawMessage, message); err != nil {
		return err
	}
	if fromController(message) {
		storageAdapter.Write(getApplicationFromControllerMessage(message), buildControllerLogMessage(message))
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
	return controllerRegex.FindStringSubmatch(message.Log)[3]
}

func buildControllerLogMessage(message *Message) string {
	l := controllerRegex.FindStringSubmatch(message.Log)
	return fmt.Sprintf("%s deis[controller]: %s %s",
		message.Time.Format(timeFormat),
		l[1],
		strings.Trim(l[4], " "))
}

func buildApplicationLogMessage(message *Message) string {
	p := podRegex.FindStringSubmatch(message.Kubernetes.PodName)
	tag := fmt.Sprintf(
		"%s.%s",
		message.Kubernetes.Labels["type"],
		message.Kubernetes.Labels["version"])
	if len(p) > 0 {
		tag = fmt.Sprintf("%s.%s", tag, p[len(p)-1])
	}
	return fmt.Sprintf("%s %s[%s]: %s",
		message.Time.Format(timeFormat),
		message.Kubernetes.Labels["app"],
		tag,
		message.Log)
}
