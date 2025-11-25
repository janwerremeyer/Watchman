package container

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	docker *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	// Check if Docker is actually running
	_, err = docker.Ping(context.TODO())
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &DockerClient{
		docker: docker,
	}, nil
}

func (c *DockerClient) ListContainers() ([]Container, error) {

	domainContainers := make([]Container, 0)

	args := filters.NewArgs()

	containers, err := c.docker.ContainerList(context.Background(), container.ListOptions{All: true, Filters: args})
	if err != nil {
		return domainContainers, err
	}

	for _, ctr := range containers {
		split := strings.Split(ctr.Image, ":")
		img := split[0]
		tag := split[len(split)-1]

		c := Container{
			Id:     ctr.ID,
			Names:  ctr.Names,
			Labels: ctr.Labels,
			State:  ctr.State,
			Status: ctr.Status,
			Image:  img,
			Tag:    tag,
		}

		domainContainers = append(domainContainers, c)
	}

	return domainContainers, nil
}

func (c *DockerClient) Start(id string) error {
	slog.Info(fmt.Sprintf("Issuing start command for container %s", id))
	return c.docker.ContainerStart(context.TODO(), id, container.StartOptions{})
}

func (c *DockerClient) Stop(id string) error {
	slog.Info(fmt.Sprintf("Issuing stop command for container %s", id))
	return c.docker.ContainerStop(context.TODO(), id, container.StopOptions{})
}

func (c *DockerClient) Purge(id string) error {
	removeOpts := container.RemoveOptions{
		Force:         true, // kill if running
		RemoveVolumes: true, // also remove attached anonymous volumes
	}

	if err := c.docker.ContainerRemove(context.TODO(), id, removeOpts); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", id, err)
	}

	return nil
}

func (c *DockerClient) ListImages() ([]Image, error) {
	imagesFromDocker, err := c.docker.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return nil, err
	}

	images := make([]Image, 0)

	for _, localImage := range imagesFromDocker {
		id := localImage.ID
		tags := localImage.RepoTags

		img := Image{
			Id:   id,
			Tags: tags,
		}

		images = append(images, img)
	}

	return images, nil
}
