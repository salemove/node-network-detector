package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

var debug = flag.Bool("debug", false, "Set true to enable debug logs")
var pingTimeout = flag.Duration("ping-timeout", 1*time.Second, "Ping node timeout")
var pingInterval = flag.Duration("ping-interval", 3*time.Second, "Ping node interval")
var nodeFetchInterval = flag.Duration("node-fetch-interval", 15*time.Second, "Nodes list fetching interval")

func main() {
	flag.Parse()
	setupLogger()

	log.Infof("Starting node-network-detector, debug: %t, pingTimeout: %v, pingInterval: %v, nodeFetchInterval: %v",
		*debug, *pingTimeout, *pingInterval, *nodeFetchInterval)

	client, err := InitKubeClient()
	if err != nil {
		panic(err.Error())
	}

	pinger := NewNodePinger(*pingTimeout, *pingInterval)

	go MonitorNodes(client, pinger, *nodeFetchInterval)
	pinger.Start()
}

func setupLogger() {
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.JSONFormatter{})
}
