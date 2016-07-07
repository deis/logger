package log

// Example log message JSON:
//
// {"log"=>"2016/05/31 01:34:43 10.164.1.1 GET / - 5074209722772702441\n", "stream"=>"stderr",
// "docker"=>{"container_id"=>"6a9069435788a05531ee2b9afbcdc73a22018af595f3203cb67e06f50103bf5f"},
// "kubernetes"=>{"namespace_name"=>"foo", "pod_id"=>"34ebc234-2423-11e6-94aa-42010a800021",
// "pod_name"=>"foo-v2-web-2ggow", "container_name"=>"foo-web", "labels"=>{"app"=>"foo",
// "heritage"=>"deis", "type"=>"web", "version"=>"v2"},
// "host"=>"gke-jchauncey-default-pool-7ae1c279-10ye"}}

// Message fields
type Message struct {
	Log        string     `json:"log"`
	Stream     string     `json:"stream"`
	Kubernetes Kubernetes `json:"kubernetes"`
	Docker     Docker     `json:"docker"`
}

// Kubernetes specific log message fields
type Kubernetes struct {
	Namespace     string            `json:"namespace_name"`
	PodID         string            `json:"pod_id"`
	PodName       string            `json:"pod_name"`
	ContainerName string            `json:"container_name"`
	Labels        map[string]string `json:"labels"`
	Host          string            `json:"host"`
}

// Docker specific log message fields
type Docker struct {
	ContainerID string `json:"container_id"`
}
