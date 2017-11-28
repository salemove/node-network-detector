package main

import (
	"flag"
	"time"

	statsd "github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
)

var debug = flag.Bool("debug", false, "Set true to enable debug logs")
var pingTimeout = flag.Duration("ping-timeout", 1*time.Second, "Ping node timeout")
var pingInterval = flag.Duration("ping-interval", 3*time.Second, "Ping node interval")
var nodeFetchInterval = flag.Duration("node-fetch-interval", 15*time.Second, "Nodes list fetching interval")
var statsdAddress = flag.String("statsd-address", "127.0.0.1:8125", "StatsD address")

func main() {
	flag.Parse()
	setupLogger()

	log.Infof("Starting node-network-detector, debug: %t, pingTimeout: %v, pingInterval: %v, nodeFetchInterval: %v",
		*debug, *pingTimeout, *pingInterval, *nodeFetchInterval)

	client, err := InitKubeClient()
	if err != nil {
		panic(err.Error())
	}

	statsd, err := statsd.New(*statsdAddress)
	if err != nil {
		panic(err)
	}

	pinger := NewNodePinger(*pingTimeout, *pingInterval, *statsd)

	go MonitorNodes(client, pinger, *nodeFetchInterval)
	pinger.Start()
}

func setupLogger() {
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	formatter := &log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "@timestamp",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
		},
	}
	log.SetFormatter(formatter)
}
