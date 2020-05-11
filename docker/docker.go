package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
)

func NewDockerLogReader(containerName string) (io.ReadCloser, error) {
	c, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	reader, err := c.ContainerLogs(ctx, containerName, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
		Tail:       "0",
	})
	if err != nil {
		return nil, err
	}
	return reader, nil
}
