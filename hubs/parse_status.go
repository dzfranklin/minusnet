package hubs

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseStatus(text string) (Status, error) {
	mapping := make(map[string]string)
	for _, line := range strings.Split(text, "\n") {
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		mapping[key] = value
	}

	return Status{
		ConnectionStatus:   mapping["Connection statuses"],
		FirmwareVersion:    mapping["Firmware version"],
		DownstreamSyncMbps: parseMbps(mapping["Downstream sync speed"]),
		UpstreamSyncMbps:   parseMbps(mapping["Upstream sync speed"]),
		NetworkUptime:      parseDuration(mapping["Network uptime"]),
		SystemUptime:       parseDuration(mapping["System uptime"]),
		DefaultGateway:     mapping["Default gateway"],
	}, nil
}

var mbpsRe = regexp.MustCompile(`^(\d+(?:\.\d+)?) Mbps$`)

func parseMbps(text string) float64 {
	m := mbpsRe.FindStringSubmatch(text)
	if m == nil {
		return -1
	}
	n, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return -1
	}
	return n
}

var durationRe = regexp.MustCompile(`^(\d+) days (\d+) Hours (\d+) Mins (\d+) Secs$`)

func parseDuration(text string) time.Duration {
	m := durationRe.FindStringSubmatch(text)
	if m == nil {
		return time.Duration(-1)
	}

	out := time.Duration(0)

	days, err := strconv.Atoi(m[1])
	if err != nil {
		return time.Duration(-1)
	}
	out += time.Duration(days) * 24 * time.Hour

	hours, err := strconv.Atoi(m[2])
	if err != nil {
		return time.Duration(-1)
	}
	out += time.Duration(hours) * time.Hour

	minutes, err := strconv.Atoi(m[3])
	if err != nil {
		return time.Duration(-1)
	}
	out += time.Duration(minutes) * time.Minute

	seconds, err := strconv.Atoi(m[4])
	if err != nil {
		return time.Duration(-1)
	}
	out += time.Duration(seconds) * time.Second

	return out
}
