package logs

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

//{"log"=>"2016/05/31 01:34:43 10.164.1.1 GET / - 5074209722772702441\n", "stream"=>"stderr", "docker"=>{"container_id"=>"6a9069435788a05531ee2b9afbcdc73a22018af595f3203cb67e06f50103bf5f"}, "kubernetes"=>{"namespace_name"=>"foo", "pod_id"=>"34ebc234-2423-11e6-94aa-42010a800021", "pod_name"=>"foo-v2-web-2ggow", "container_name"=>"foo-web", "labels"=>{"app"=>"foo", "heritage"=>"deis", "type"=>"web", "version"=>"v2"}, "host"=>"gke-jchauncey-default-pool-7ae1c279-10ye"}}

// Message json log message from fluentd
type Message struct {
	Log        string     `json:"log"`
	Stream     string     `json:"stream"`
	Kubernetes Kubernetes `json:"kubernetes"`
	Docker     Docker     `json:"docker"`
}

// Kuberentes decorated log message
type Kubernetes struct {
	Namespace     string            `json:"namespace_name"`
	PodID         string            `json:"pod_id"`
	PodName       string            `json:"pod_name"`
	ContainerName string            `json:"container_name"`
	Labels        map[string]string `json:"labels"`
	Host          string            `json:"host"`
}

// Docker decorated log messaage
type Docker struct {
	ContainerID string `json:"container_id"`
}

// Logger holds the storage adapter we are using
type Logger struct {
	StorageAdapter storage.Adapter
}

// NewLogger returns a Logger object with the storage adapter set
func NewLogger(storageAdapter storage.Adapter) (*Logger, error) {
	if storageAdapter == nil {
		return nil, fmt.Errorf("No storage adapter specified.")
	}
	return &Logger{
		StorageAdapter: storageAdapter,
	}, nil
}

// ReadLogs returns a specified number of log lines (if available) for a specified app by
// delegating to the server's underlying storage.Adapter.
func (logger *Logger) ReadLogs(app string, lines int) ([]string, error) {
	return logger.StorageAdapter.Read(app, lines)
}

// DestroyLogs deletes all logs for a specified app by delegating to the server's underlying
// storage.Adapter.
func (logger *Logger) DestroyLogs(app string) error {
	return logger.StorageAdapter.Destroy(app)
}

// WriteLog - takes string which is JSON and stores it in the storage adapter.
func (logger *Logger) WriteLog(rawMessage []byte) error {
	message := new(Message)
	if err := json.Unmarshal(rawMessage, message); err != nil {
		return err
	}
	if fromController(message) {
		logger.StorageAdapter.Write(getApplicationFromControllerMessage(message), message.Log)
	} else {
		labels := message.Kubernetes.Labels
		logger.StorageAdapter.Write(labels["app"], buildApplicationLogMessage(message))
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
