package syslog

import (
	"fmt"
	"gopkg.in/mcuadros/go-syslog.v2"
	"net/url"
)

func openListener(s *syslog.Server, c string) error {
	u, err := url.Parse(c)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "tcp":
		err := s.ListenTCP(u.Host)
		if err != nil {
			return err
		}

	case "udp":
		err := s.ListenUDP(u.Host)
		if err != nil {
			return err
		}

	case "unix":
		return fmt.Errorf("Not implemented")

	default:
		return fmt.Errorf("syslog server should be in format unix/tcp/udp://127.0.0.1:5533")
	}

	return nil
}

func SyslogRunner(conn string) (syslog.LogPartsChannel, error) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()

	//RFC3164 or RFC5424 or RFC6587. nginx works on RFC3164
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)

	err := openListener(server, conn)
	if err != nil {
		return nil, err
	}

	err = server.Boot()
	if err != nil {
		return nil, err
	}

	return channel, nil
}
