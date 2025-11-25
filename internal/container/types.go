package container

type Client interface {
	ListContainers() ([]Container, error)
	Start(id string) error
	Stop(id string) error
	Purge(id string) error
	ListImages() ([]Image, error)
}

type Registry interface {
	ListTags(image string) ([]string, error)
	PullImage(imageWithTag string) error
}

type Container struct {
	Id     string            `json:"id"`
	Names  []string          `json:"names"`
	Labels map[string]string `json:"labels"`
	State  string            `json:"state"`
	Status string            `json:"status"`
	Image  string            `json:"image"`
	Tag    string            `json:"tag"`
}

type Image struct {
	Id   string
	Tags []string
}
