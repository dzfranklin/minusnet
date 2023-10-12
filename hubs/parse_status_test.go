package hubs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const exampleStatusText = `Connection status,Connected
Connection type,Fibre Broadband (VDSL)
Firmware version,v0.10.00.04201-PN (Thu Apr 20 16:54:19 2023)
Serial number,+108417+2216005167
Downstream sync speed,45.3 Mbps
Upstream sync speed,9.3 Mbps
Network uptime,0 days 1 Hours 41 Mins 53 Secs
System uptime,0 days 1 Hours 43 Mins 20 Secs
Broadband IP address,146.198.207.137
Default gateway,172.16.15.122
Primary DNS,212.159.6.9
Secondary DNS,212.159.6.10
`

func TestParseExampleStatusText(t *testing.T) {
	want := Status{
		ConnectionStatus:   "Connected",
		FirmwareVersion:    "v0.10.00.04201-PN (Thu Apr 20 16:54:19 2023)",
		DownstreamSyncMbps: 45.3,
		UpstreamSyncMbps:   9.3,
		NetworkUptime:      1*time.Hour + 41*time.Minute + 53*time.Second,
		SystemUptime:       1*time.Hour + 43*time.Minute + 20*time.Second,
		DefaultGateway:     "172.16.15.122",
	}
	got, err := parseStatus(exampleStatusText)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
