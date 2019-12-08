package syslog

import (
	"fmt"
	"net/url"

	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
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

// Listen opens up a new syslog server on either a TCP or UDP port
func Listen(conn string, formatSpec string) (syslog.LogPartsChannel, *syslog.Server, error) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()

	var format format.Format = syslog.Automatic

	switch formatSpec {
	case "rfc3164":
		format = syslog.RFC3164
	case "rfc5424":
		format = syslog.RFC5424
	case "rfc6587":
		format = syslog.RFC6587
	case "auto":
		format = syslog.Automatic
	case "":
		format = syslog.Automatic
	default:
		return nil, nil, fmt.Errorf("unknown syslog format: '%s'", format)
	}

	//RFC3164 or RFC5424 or RFC6587. nginx works on RFC3164
	server.SetFormat(format)
	server.SetHandler(handler)

	err := openListener(server, conn)
	if err != nil {
		return nil, nil, err
	}

	err = server.Boot()
	if err != nil {
		return nil, nil, err
	}

	return channel, server, nil
}
