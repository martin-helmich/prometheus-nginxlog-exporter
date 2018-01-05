package tail

import (
	"github.com/hpcloud/tail"
)

type Follower interface {
	Lines() chan *tail.Line
	OnError(func (error))
}

type followerImpl struct {
	filename string
	t *tail.Tail
}

func NewFollower(filename string) (Follower, error) {
	f := &followerImpl{
		filename: filename,
	}

	if err := f.start(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *followerImpl) start() error {
	t, err := tail.TailFile(f.filename, tail.Config{
		Follow: true,
		ReOpen: true,
		Poll:   true,
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

func (f *followerImpl) Lines() chan *tail.Line {
	return f.t.Lines
}