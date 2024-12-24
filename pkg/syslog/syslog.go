package syslog

import (
	"fmt"
	"net/url"
	"os"

	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func openListener(s *syslog.Server, c string) (func() error, error) {
	u, err := url.Parse(c)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "tcp":
		return nil, s.ListenTCP(u.Host)

	case "udp":
		return nil, s.ListenUDP(u.Host)

	case "unix":
		socketPath := u.Host + u.Path

		if err := s.ListenUnixgram(socketPath); err != nil {
			return nil, err
		}

		return func() error {
			return os.Remove(socketPath)
		}, nil

	default:
		return nil, fmt.Errorf("syslog server should be in format unix/tcp/udp://127.0.0.1:5533")
	}
}

// Listen opens up a new syslog server on either a TCP or UDP port
func Listen(conn string, formatSpec string) (syslog.LogPartsChannel, *syslog.Server, func() error, error) {
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
		return nil, nil, nil, fmt.Errorf("unknown syslog format: '%s'", format)
	}

	//RFC3164 or RFC5424 or RFC6587. nginx works on RFC3164
	server.SetFormat(format)
	server.SetHandler(handler)

	closeListener, err := openListener(server, conn)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = server.Boot(); err != nil {
		return nil, nil, nil, err
	}

	stopFn := func() error {
		if err := server.Kill(); err != nil {
			return fmt.Errorf("failed to kill syslog server: %w", err)
		}

		if closeListener != nil {
			return closeListener()
		}

		return nil
	}

	return channel, server, stopFn, nil
}
