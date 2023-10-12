package main

import (
	"fmt"
	"github.com/dzfranklin/minusnet/hubs"
	"github.com/dzfranklin/minusnet/stores"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"os"
)

var hub *hubs.Hub
var store *stores.Store

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	l := zap.S()

	hub = &hubs.Hub{
		Addr:   "192.168.1.254",
		Serial: "+108417+2216005167",
	}

	dataDir := os.Getenv("MINUSNET_DATA_DIR")
	if dataDir == "" {
		dataDir = "/tmp/minusnet"
	}
	store, err = stores.New(dataDir)
	defer func() {
		if err := store.Close(); err != nil {
			l.Errorf("failed to close store: %s", err)
		}
	}()

	if len(os.Args) < 2 {
		println("usage: minusnet <command>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "scrape":
		scrapeMain()
	default:
		fmt.Printf("unknown command: %s", os.Args[1])
		os.Exit(1)
	}
}

func scrapeMain() {
	l := zap.S()
	l.Info("scraping")

	status, err := hub.RequestStatus()
	if err != nil {
		l.Fatalf("failed to request status: %s", err)
	}
	l.Infof("got status:\n%s", status)

	if err := store.RecordStatus(status); err != nil {
		l.Fatalf("failed to insert status: %s", err)
	}
	l.Info("stored status")
}
