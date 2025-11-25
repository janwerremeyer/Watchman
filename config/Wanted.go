package config

type Wanted struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Tag         string            `yaml:"tag"`
	AutoUpdate  bool              `yaml:"auto_update"`
	Replicas    int               `yaml:"replicas"`
	Environment map[string]string `yaml:"environment"`
}
