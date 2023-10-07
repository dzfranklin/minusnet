package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path"
	"time"

	"github.com/dzfranklin/minusnet/status"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	l := zap.S()

	if len(os.Args) < 2 {
		l.Fatalf("usage: %s <command>", os.Args[0])
	}
	cmd := os.Args[1]

	dataDir := os.Getenv("MINUSNET_DATA_DIR")
	if dataDir == "" {
		dataDir = "/tmp/minusnet"
	}
	l.Infof("using data dir: %s", dataDir)
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		l.Fatalf("failed to create data dir (%s): %s", dataDir, err)
	}

	db, err := sql.Open("sqlite3", path.Join(dataDir, "minusnet.db"))
	if err != nil {
		l.Fatalf("failed to open db: %s", err)
	}
	defer db.Close()

	err = initDb(db)
	if err != nil {
		l.Fatalf("failed to initialize db: %s", err)
	}

	switch cmd {
	case "scrape":
		scrapeMain(l, db)
	default:
		l.Fatalf("unknown command: %s", cmd)
	}
}

func initDb(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS status (
			id INTEGER PRIMARY KEY,
			connection_status TEXT,
			firmware_version TEXT,
			downstream_sync_mbps REAL,
			upstream_sync_mbps REAL,
			network_uptime_secs INTEGER,
			system_uptime_secs INTEGER,
			default_gateway TEXT,
			created_at INTEGER
		)`)
	if err != nil {
		return fmt.Errorf("create status table: %w", err)
	}
	return nil
}

func scrapeMain(l *zap.SugaredLogger, db *sql.DB) {
	l.Info("scraping")

	status, err := status.Request()
	if err != nil {
		l.Fatalf("failed to request status: %s", err)
	}
	l.Infof("status:\n%s", status)

	_, err = db.Exec(`
		INSERT INTO status (
			connection_status,
			firmware_version,
			downstream_sync_mbps,
			upstream_sync_mbps,
			network_uptime_secs,
			system_uptime_secs,
			default_gateway,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		status.ConnectionStatus,
		status.FirmwareVersion,
		status.DownstreamSyncMbps,
		status.UpstreamSyncMbps,
		math.Round(status.NetworkUptime.Seconds()),
		math.Round(status.SystemUptime.Seconds()),
		status.DefaultGateway,
		time.Now().Unix(),
	)
	if err != nil {
		l.Fatalf("failed to insert status: %s", err)
	}
	l.Info("wrote status to db")
}
