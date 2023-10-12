package stores

import (
	"database/sql"
	"fmt"
	"github.com/dzfranklin/minusnet/hubs"
	"go.uber.org/zap"
	"math"
	"os"
	"path"
	"time"
)

type Store struct {
	db *sql.DB
}

func New(dataDir string) (*Store, error) {
	if dataDir == "" {
		panic("dataDir is required")
	}
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("create dataDir (%s): %w", dataDir, err)
	}
	fname := path.Join(dataDir, "minusnet.db")
	db, err := sql.Open("sqlite3", fname)
	if err != nil {
		return nil, fmt.Errorf("open db (%s): %w", fname, err)
	}
	store := &Store{db}
	if err = store.init(); err != nil {
		return nil, fmt.Errorf("init db (%s): %w", fname, err)
	}
	zap.S().Infof("opened %s", fname)
	return store, nil
}

func (s *Store) init() error {
	_, err := s.db.Exec(`
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
		return fmt.Errorf("status table: %w", err)
	}
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS restarts (
			id INTEGER PRIMARY KEY,
			created_at INTEGER
		)`)
	if err != nil {
		return fmt.Errorf("restarts table: %w", err)
	}
	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) RecordStatus(status hubs.Status) error {
	_, err := s.db.Exec(`
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
		return fmt.Errorf("insert status: %w", err)
	}
	return nil
}

func (s *Store) RecordRestart() error {
	_, err := s.db.Exec(`INSERT INTO restarts (created_at) VALUES (?)`, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("insert restart: %w", err)
	}
	return nil
}
