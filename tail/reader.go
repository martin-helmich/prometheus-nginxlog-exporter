package tail

import (
	"io"
	"strings"
)

type readerFollower struct {
	line          chan string
	logReader     io.Reader
	bufferSize    int
	errorHandlers []func(error)
}

// NewReaderFollower builds a new reader follower from a reader that reads log data
func NewReaderFollower(logReader io.ReadCloser) (Follower, error) {
	s := &readerFollower{
		line:          make(chan string),
		logReader:     logReader,
		bufferSize:    32 * 1024,
		errorHandlers: make([]func(error), 0),
	}
	return s, nil
}

func (d *readerFollower) OnError(cb func(error)) {
	d.errorHandlers = append(d.errorHandlers, cb)
}

func (d *readerFollower) Lines() chan string {
	go func() {
		lines := []string{""}
		for {
			buf := make([]byte, d.bufferSize)
			numRead, err := d.logReader.Read(buf)
			if err != nil {
				for _, handler := range d.errorHandlers {
					handler(err)
				}
			}
			s := string(buf[:numRead])
			newLines := strings.Split(s, "\n")
			if len(newLines) > 0 {
				lines[len(lines)-1] = lines[len(lines)-1] + newLines[0]
				lines = append(lines, newLines[1:]...)
			}

			for len(lines) > 1 {
				d.line <- lines[0]
				lines = lines[1:]
			}
		}
	}()
	return d.line
}
