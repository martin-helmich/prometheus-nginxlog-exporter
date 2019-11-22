package tail

import (
	"os"

	"github.com/hpcloud/tail"
)

type followerImpl struct {
	filename string
	t        *tail.Tail
	line     chan string
}

// NewFollower creates a new Follower instance for a given file (given by name)
func NewFileFollower(filename string) (Follower, error) {
	f := &followerImpl{
		filename: filename,
		line:     make(chan string),
	}

	if err := f.start(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *followerImpl) start() error {
	var seekInfo *tail.SeekInfo

	_, err := os.Stat(f.filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		seekInfo = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	}

	t, err := tail.TailFile(f.filename, tail.Config{
		Follow:   true,
		ReOpen:   true,
		Poll:     true,
		Location: seekInfo,
	})

	if err != nil {
		return err
	}

	f.t = t
	return nil
}

func (f *followerImpl) OnError(cb func(error)) {
	go func() {
		err := f.t.Wait()
		if err != nil {
			cb(err)
		}
	}()
}

func (f *followerImpl) Lines() chan string {
	go func() {
		for n := range f.t.Lines {
			f.line <- n.Text
		}
	}()
	return f.line
}
