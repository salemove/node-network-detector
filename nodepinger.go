package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	// Using bastjan fork because sparrc has not applied
	// https://github.com/sparrc/go-ping/pull/15 and other useful fixes yet
	statsd "github.com/DataDog/datadog-go/statsd"
	ping "github.com/bastjan/go-ping"
)

type NodePinger struct {
	Nodes        []string
	PingTimeout  time.Duration
	PingInterval time.Duration
	StatsD       statsd.Client
}

func NewNodePinger(pingTimeout time.Duration, pingInterval time.Duration, statsdClient statsd.Client) *NodePinger {
	pinger := NodePinger{
		Nodes:        make([]string, 0),
		PingTimeout:  pingTimeout,
		PingInterval: pingInterval,
		StatsD:       statsdClient,
	}
	return &pinger
}

func (p *NodePinger) Start() {
	log.Info("NodePinger started")

	for {
		var wg sync.WaitGroup
		wg.Add(len(p.Nodes))

		for _, ip := range p.Nodes {
			go p.pingNode(&wg, ip)
		}

		wg.Wait()

		time.Sleep(p.PingInterval)
	}
}

func (p *NodePinger) pingNode(wg *sync.WaitGroup, ip string) {
	defer wg.Done()

	log.Debugf("Pinging node %s", ip)

	ctx := context.TODO()
	pinger, err := ping.NewPinger(ctx, ip)
	if err != nil {
		log.Warnf("Unable to create pinger for %s, %s", ip, err)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = p.PingTimeout
	pinger.Run()

	stats := pinger.Statistics()
	log.Debugf("Pinging node %s results: %+v", ip, stats)

	if stats.PacketLoss > 0 || stats.PacketsRecv == 0 || stats.PacketsSent == 0 {
		log.Warnf("Unable to reach %s, %+v", ip, stats)
	} else {
		if stats.PacketsRecv > 1 {
			log.Warnf("Got more packets that expected %s, %+v", ip, stats)
		}

		tags := []string{
			"target:" + ip,
		}
		p.StatsD.Timing("sm.node-network-detector.ping.duration", stats.AvgRtt, tags, 1)
	}
}
