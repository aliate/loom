package fluentd

type DockerLogEntry struct {
	Hostname      string `json:"hostname"`
	ContainerTag  string `json:"container_tag"`
	ContainerId   string `json:"container_id"`
	ContainerName string `json:"container_name"`
	Source        string `json:"source"`
	Timestamp     int64  `json:"timestamp"`
	Log           string `json:"log"`
}
