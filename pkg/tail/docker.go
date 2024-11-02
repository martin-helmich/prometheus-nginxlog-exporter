package tail

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
	"io"
	"strings"
)

// Bytes written to channelLineWriter are split into lines of strings and sent to a channel
type channelLineWriter struct {
	c    chan string
	last string
}

func (w *channelLineWriter) Write(p []byte) (int, error) {
	// Split log data into lines
	// Send all but the last line to the channel
	// Concatenate the last line from the previous write to the first line from the current write
	lines := strings.Split(string(p), "\n")
	if len(lines) > 0 {
		lines[0] = w.last + lines[0]

		// send all but last line and remove sent lines
		for len(lines) > 1 {
			w.c <- lines[0]
			lines = lines[1:]
		}

		w.last = lines[0]
	}

	return len(p), nil
}

type dockerFollower struct {
	container     string
	line          chan string
	logReader     io.ReadCloser
	errorHandlers []func(error)
	tty           bool
}

// NewDockerFollower builds a new follower that reads data from dockers default logging driver
func NewDockerFollower(containerName string) (Follower, error) {
	c, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	ins, err := c.ContainerInspect(ctx, containerName)
	if err != nil {
		return nil, err
	}

	reader, err := c.ContainerLogs(ctx, containerName, container.LogsOptions{
		ShowStderr: false,
		ShowStdout: true,
		Follow:     true,
		Tail:       "0",
	})
	if err != nil {
		return nil, err
	}

	s := &dockerFollower{
		container:     containerName,
		line:          make(chan string),
		logReader:     reader,
		errorHandlers: make([]func(error), 0),
		tty:           ins.Config.Tty,
	}
	return s, nil
}

func (d *dockerFollower) OnError(cb func(error)) {
	d.errorHandlers = append(d.errorHandlers, cb)
}

func (d *dockerFollower) Lines() chan string {
	w := &channelLineWriter{
		c:    d.line,
		last: "",
	}
	go func() {
		for {
			// Read logs
			// If the docker container is using a TTY there is a single output stream we can read directly
			// If it is not using a TTY we need to use the stdcopy.StdCopy to handle the multiplexed data
			// See https://godoc.org/github.com/docker/docker/client#Client.ContainerLogs for more info
			var err error
			if d.tty {
				_, err = io.Copy(w, d.logReader)
			} else {
				_, err = stdcopy.StdCopy(w, w, d.logReader)
			}
			if err == io.EOF {
				err = errors.New(fmt.Sprintf("end of logs for container %s", d.container))
			}
			if err != nil {
				for _, handler := range d.errorHandlers {
					handler(err)
				}
			}
		}
	}()
	return d.line
}
