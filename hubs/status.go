package hubs

import (
	"fmt"
	"strings"
	"time"
)

type Status struct {
	ConnectionStatus   string
	FirmwareVersion    string
	DownstreamSyncMbps float64
	UpstreamSyncMbps   float64
	NetworkUptime      time.Duration
	SystemUptime       time.Duration
	DefaultGateway     string
}

func (s Status) String() string {
	return strings.Join([]string{
		"Connection status: " + s.ConnectionStatus,
		"Firmware version: " + s.FirmwareVersion,
		fmt.Sprintf("Downstream sync speed: %.1f Mbps", s.DownstreamSyncMbps),
		fmt.Sprintf("Upstream sync speed: %.1f Mbps", s.UpstreamSyncMbps),
		fmt.Sprintf("Network uptime: %s", s.NetworkUptime),
		fmt.Sprintf("System uptime: %s", s.SystemUptime),
		"Default gateway: " + s.DefaultGateway,
	}, "\n")
}
