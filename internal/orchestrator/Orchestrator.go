package orchestrator

import (
	"log/slog"
	"watchman/config"
	"watchman/internal/container"
	"watchman/internal/gofy"
)

type Orchestrator struct {
	client   container.Client
	registry container.Registry
	wanted   []config.Wanted
}

func NewOrchestrator(client container.Client, registry container.Registry, wanted []config.Wanted) *Orchestrator {
	return &Orchestrator{
		client:   client,
		registry: registry,
		wanted:   wanted,
	}
}

func (o *Orchestrator) Run() {
	existing, err := o.client.ListContainers()
	if err != nil {
		slog.Error(err.Error())
	}

	for _, wanted := range o.wanted {
		found := make([]container.Container, 0)
		for _, cnt := range existing {
			if cnt.Image == wanted.Image && cnt.Labels["watchman_name"] == wanted.Name {
				found = append(found, cnt)
			}
		}

		validStates := []string{"created", "running", "restarting"}

		for _, cnt := range found {
			if !gofy.Contains(cnt.State, validStates) {
				_ = o.client.Purge(cnt.Id)
			}
		}

		healthy := 0
		for _, cnt := range found {
			if gofy.Contains(cnt.State, validStates) {
				healthy += 1
			}
		}

		if healthy >= wanted.Replicas {
			//Purge overhead
		}

		if healthy <= wanted.Replicas {
			//Spawn new
		}
	}

}
