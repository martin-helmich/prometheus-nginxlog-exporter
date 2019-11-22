package tail

import (
	"gopkg.in/mcuadros/go-syslog.v2"
)

type syslogFollower struct {
	tag  string
	line chan string

	channel syslog.LogPartsChannel
	server  syslog.Server
}

func NewSyslogFollower(tag string, channel syslog.LogPartsChannel) (Follower, error) {
	s := &syslogFollower{
		tag:     tag,
		channel: channel,
		line:    make(chan string),
	}
	return s, nil
}

func (s *syslogFollower) OnError(cb func(error)) {
	go func() {
		err := s.server.GetLastError()
		if err != nil {
			cb(err)
		}
	}()
}

func (s *syslogFollower) Lines() chan string {
	go func() {
		for line := range s.channel {
			if _, ok := line["tag"].(string); !ok {
				continue
			}

			if line["tag"].(string) == s.tag {
				s.line <- line["content"].(string)
			}
		}
	}()
	return s.line
}
